import {
  AfterViewInit,
  Component,
  ElementRef,
  EventEmitter,
  OnInit,
  Output,
  Renderer2,
  ViewChild
} from '@angular/core';
import { Filter } from '../../../../lib/filter-dsl/filter';
import { PrintPlainText } from '../../../../lib/filter-dsl/printer';
import { CSTNode, NodeKind, ParseCST, ParseError, renderCST, ValueNode } from '../../../../lib/filter-dsl/cst';
import { Tag } from '../../../../lib/filter-dsl/scanner';
import { ActivatedRoute } from '@angular/router';
import { Observable } from 'rxjs';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { map } from 'rxjs/operators';
import { ValueKind } from '../../../../lib/filter-dsl/value';

@Component({
  selector: 'app-dsl-search',
  templateUrl: './dsl-search.component.html',
  styleUrls: ['./dsl-search.component.scss'],
})
export class DslSearchComponent implements OnInit, AfterViewInit {

  @Output()
  executeQuery: EventEmitter<string> = new EventEmitter<string>();

  @ViewChild('editableContent')
  editableContent: ElementRef;

  keyboardEvents: EventEmitter<KeyboardEvent> = new EventEmitter();

  initialQuery: string;

  cst: CSTNode = null;
  filter: Filter = null;
  error: string = null;
  info: string = null;
  inputActive = false;
  inputEmpty = true;

  sampleQueries: string[] = [
    `actor = "ricky" and content ~= "chimpanzee that"`,
    `actor = "steve" and content = "arbitrary"`,
    `type = "song" and content = "feeder"`
  ];

  private caretPos: number;
  private activeValueNode: ValueNode;

  // data needed for value dropdowns
  dropdownActive: boolean = false;
  dropdownFieldName: string;
  dropdownFilter: string;
  dropdownValueSource: (fieldName: string, prefix: string) => Observable<string[]>;
  dropdownFieldWhitelist: string[] = ['actor', 'type'];

  constructor(private renderer: Renderer2, route: ActivatedRoute, private apiClient: SearchAPIClient) {
    route.queryParamMap.subscribe((params) => {
      if (params.get('q') === null) {
        return;
      }
      this.initialQuery = params.get('q');
    });
  }

  ngOnInit(): void {
  }

  ngAfterViewInit(): void {
    this.parse();
  }

  activate() {
    this.inputActive = true;
    this.inputEmpty = this.editableContent.nativeElement.innerText.trim().length === 0;
    this.parse();
    this.updateDropdown();
  }

  deactivate() {
    this.inputActive = false;
    this.resetDropdown();
  }

  emitQuery() {
    if (!this.cst) {
      return;
    }
    this.executeQuery.emit(this.cst.string());
  }

  parse() {
    if (this.editableContent.nativeElement.innerText.length === 0) {
      return;
    }
    this.caretPos = this.getCaretPosition(this.editableContent.nativeElement);
    try {
      this.cst = ParseCST(this.editableContent.nativeElement.innerText);
      if (this.cst != null) {
        this.editableContent.nativeElement.innerHTML = '';
        this.renderer.appendChild(this.editableContent.nativeElement, renderCST(this.renderer, this.cst));
        this.clearNotices();
        this.moveCaretTo(this.caretPos);
      }
    } catch (e) {
      if (e instanceof ParseError) {
        if (e.cause.tag === Tag.EOF) {
          this.setWarning(`incomplete query (${e.reason})`);
        } else {
          this.setError(`${e.reason} pos: ${e.cause.start} near: "${e.cause.lexeme || e.cause.tag}"`);
        }
      } else {
        this.setError(e.reason);
      }
      this.cst = null;
    }
  }

  clearNotices() {
    this.error = null;
    this.info = null;
  }

  setError(msg: string) {
    this.error = msg;
  }

  setWarning(msg: string) {
    this.clearNotices();
    this.info = msg;
  }

  filterText(): string {
    return PrintPlainText(this.filter);
  }

  onKeydown(key: KeyboardEvent): boolean {

    this.resetDropdown();

    switch (key.code) {
      case 'ArrowDown':
        if (!this.dropdownActive) {
          this.activate();
        }
        return false;
      case 'ArrowUp':
        if (!this.dropdownActive) {
          this.activate();
        }
        return false;
      case 'Enter':
        this.activate();
        if (!this.dropdownActive) {
          this.emitQuery();
        }
        return false;
      case 'Escape':
        this.deactivate();
        break;
      default:
        this.activate();
    }
  }

  onKeypress(key: KeyboardEvent) {

    // forward event to any other components e.g. value-list
    this.keyboardEvents.next(key);

    if (['Shift', 'Alt', 'Control', 'Escape'].indexOf(key.code) !== -1) {
      return;
    }
    this.activate();
  }

  private moveCaretTo(position: number): void {

    // move to end (no idea why it doesn't work in the normal way)
    if (position >= this.editableContent.nativeElement.textContent.length) {
      let range, selection;
      range = document.createRange();
      range.selectNodeContents(this.editableContent.nativeElement);
      range.collapse(false);
      selection = window.getSelection();
      selection.removeAllRanges();
      selection.addRange(range);
      return;
    }

    const node = this.getTextNodeAtPosition(
      this.editableContent.nativeElement,
      position
    );
    const sel = window.getSelection();
    sel.collapse(node.node, node.position);
  }

  private getCaretPosition(container: HTMLElement): number {
    const selection = window.getSelection();
    if (selection.rangeCount === 0) {
      return 0;
    }
    const range = selection.getRangeAt(0);
    const selected = range.toString().length;
    const preCaretRange = range.cloneRange();
    preCaretRange.selectNodeContents(container);
    preCaretRange.setEnd(range.endContainer, range.endOffset);

    if (selected) {
      return preCaretRange.toString().length - selected;
    }
    return preCaretRange.toString().length;
  }

  private getTextNodeAtPosition(root, index) {
    let lastNode = null;
    const treeWalker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT, {
      acceptNode: (elem: Node): number => {
        if (index >= elem.textContent.length) {
          index -= elem.textContent.length;
          lastNode = elem;
          return NodeFilter.FILTER_REJECT;
        }
        return NodeFilter.FILTER_ACCEPT;
      }
    });
    const c = treeWalker.nextNode();
    return { node: (c ? c : root), position: (c ? index : 0) };
  }

  findActiveValue(): ValueNode {
    this.activeValueNode = null;
    if (!this.cst) {
      return null;
    }
    this.cst.walk((v) => {
      if (this.caretPos > v.startPos() && this.caretPos <= v.endPos()) {
        if (v.kind === NodeKind.Value) {
          this.activeValueNode = v as ValueNode;
        }
      }
    });
    return this.activeValueNode;
  }

  applyQuery(q: string) {
    this.editableContent.nativeElement.textContent = q;
    this.parse();
  }

  resetDropdown() {
    this.dropdownFieldName = null;
    this.dropdownFilter = null;
    this.dropdownValueSource = null;
    this.dropdownActive = false;
  }

  updateDropdown() {
    const el = this.findActiveValue();
    if (!el) {
      return;
    }
    const field = el.parent.children[0].string();
    if (this.dropdownFieldWhitelist.indexOf(field) === -1) {
      return;
    }

    if (field === this.dropdownFieldName && this.dropdownFilter === el.v.v) {
      return;
    }

    this.resetDropdown();
    this.dropdownActive = true;
    this.dropdownFieldName = field;
    this.dropdownFilter = el.v.v;
    this.dropdownValueSource = (fieldName: string, prefix: string): Observable<string[]> => {
      return this.apiClient.searchServiceListFieldValues({
        field: fieldName,
        prefix: prefix
      }).pipe(map((res => res.values.map(val => val.value))));
    };
  }

  dropdownValueSelected(value: string[]) {

    if (!value || value.length === 0) {
      return;
    }

    const el = this.findActiveValue();
    if (!el) {
      return;
    }
    // ignore multi-select for now. The dsl doesn't have IN support anyway.

    let strVal: string = '';
    switch (el.v.kind) {
      case ValueKind.String:
        strVal = `"${value[0]}"`;
        break;
      default:
        strVal = value[0];
    }

    this.editableContent.nativeElement.textContent = spliceString(
      this.editableContent.nativeElement.textContent,
      el.startPos(),
      el.endPos(),
      strVal,
    );

    // move caret to end of the new token
    this.moveCaretTo(el.startPos() + strVal.length);

    // re-parse input
    this.parse();

    // clear dropdown
    this.resetDropdown();
  }
}

function spliceString(str: string, start: number, end: number, replace: string) {
  if (start < 0) {
    start = str.length + start;
    start = start < 0 ? 0 : start;
  }
  return str.slice(0, start) + (replace || '') + str.slice(end);
}
