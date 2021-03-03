import { Component, ElementRef, EventEmitter, OnInit, Output, Renderer2, ViewChild } from '@angular/core';
import { Filter } from '../../../../lib/filter-dsl/filter';
import { PrintPlainText } from '../../../../lib/filter-dsl/printer';
import { CSTNode, ParseCST, ParseError, renderCST } from '../../../../lib/filter-dsl/cst';
import { Tag } from '../../../../lib/filter-dsl/scanner';

@Component({
  selector: 'app-dsl-search',
  templateUrl: './dsl-search.component.html',
  styleUrls: ['./dsl-search.component.scss']
})
export class DslSearchComponent implements OnInit {

  @Output()
  executeQuery: EventEmitter<string> = new EventEmitter<string>();

  @ViewChild('editableContent')
  editableContent: ElementRef;

  @ViewChild('renderedFilter')
  renderedFilter: ElementRef;

  cst: CSTNode = null;
  filter: Filter = null;
  error: string = null;
  info: string = null;

  private inputActive = false;
  private caretPos: number;
  private contentLength: number;

  constructor(private renderer: Renderer2) {
  }

  ngOnInit(): void {
  }

  activate() {
    this.inputActive = true;
    this.parse();
  }

  deactivate() {
    this.inputActive = false;
  }

  emitQuery() {
    this.executeQuery.emit(null);
  }

  parse() {
    const originalCaretPos = this.getCaretPosition(this.editableContent.nativeElement);
    const originalText = this.editableContent.nativeElement.innerText;
    try {
      this.cst = ParseCST(this.editableContent.nativeElement.innerText);
      if (this.cst != null) {

        console.log(this.cst);

        this.editableContent.nativeElement.innerHTML = '';
        this.renderer.appendChild(this.editableContent.nativeElement , renderCST(this.renderer, this.cst));
        this.clearNotices();
        this.moveCaretTo(originalCaretPos);
      }
    } catch (e) {
      console.error(e);
      if (e instanceof ParseError) {
        if (e.cause.tag === Tag.EOF) {
          this.setWarning(`incomplete query (${e.reason})`);
        } else {
          this.setError(`${e.reason} pos: ${e.cause.start} near: "${e.cause.lexeme || e.cause.tag}"`)
        }
      } else {
        this.setError(e.reason);
      }
      this.filter = null;

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

}
