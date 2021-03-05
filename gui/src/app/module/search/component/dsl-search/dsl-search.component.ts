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
import { CSTNode, ParseCST, ParseError, renderCST } from '../../../../lib/filter-dsl/cst';
import { Tag } from '../../../../lib/filter-dsl/scanner';

@Component({
  selector: 'app-dsl-search',
  templateUrl: './dsl-search.component.html',
  styleUrls: ['./dsl-search.component.scss']
})
export class DslSearchComponent implements OnInit, AfterViewInit {

  @Output()
  executeQuery: EventEmitter<string> = new EventEmitter<string>();

  @ViewChild('editableContent')
  editableContent: ElementRef;

  cst: CSTNode = null;
  filter: Filter = null;
  error: string = null;
  info: string = null;
  inputActive = false;

  sampleQueries: string[] = [
    `actor = "ricky" and content ~= "chimpanzee that`,
    `actor = "steve" and content = "arbitrary"`,
    `type = "song" and content = "feeder"`
  ];

  private caretPos: number;
  private activeNode: CSTNode = null;

  constructor(private renderer: Renderer2) {
  }

  ngOnInit(): void {

  }

  ngAfterViewInit(): void {
    this.parse();
  }

  activate() {
    this.inputActive = true;
    this.parse();
  }

  deactivate() {
    this.inputActive = false;
  }

  emitQuery() {
    if (!this.cst) {
      return;
    }
    this.executeQuery.emit(this.cst.string());
  }

  parse() {
    if (this.editableContent.nativeElement.innerText.length === 0) {
      return
    }
    //this.caretPos = this.getCaretPosition(this.editableContent.nativeElement);
    try {
      this.cst = ParseCST(this.editableContent.nativeElement.innerText);
      if (this.cst != null) {
        this.editableContent.nativeElement.innerHTML = '';
        this.renderer.appendChild(this.editableContent.nativeElement, renderCST(this.renderer, this.cst));
        this.clearNotices();
        //this.moveCaretTo(this.caretPos);
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
    switch (key.code) {
      case 'ArrowDown':
      case 'ArrowUp':
      case 'Enter':
        this.activate();
        this.emitQuery();
        return false;
      case 'Escape':
        this.deactivate();
        break;
    }
  }

  onKeypress(key: KeyboardEvent) {

    // forward event to any other components e.g. value-list
    //this.keyboardEvents.next(evt);

    if (['Shift', 'Alt', 'Control'].indexOf(key.code) !== -1) {
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

  findActiveElement() {
    if (!this.cst) {
      return null;
    }
    this.activeNode = null;
    this.cst.walk((v) => {
      console.log(v.startPos(), v.endPos(), this.caretPos);
      if (this.caretPos > v.startPos() && this.caretPos <= v.endPos()) {
        this.activeNode = v;
      }
    });

    return this.activeNode;
  }

  applyQuery(q: string) {
    this.editableContent.nativeElement.textContent = q;
    this.parse();
  }

}
