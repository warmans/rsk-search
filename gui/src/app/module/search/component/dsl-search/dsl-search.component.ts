import { Component, ElementRef, EventEmitter, OnInit, Output, Renderer2, ViewChild } from '@angular/core';
import { Filter } from '../../../../lib/filter-dsl/filter';
import { Parse } from '../../../../lib/filter-dsl/parser';
import { PrintHTML, PrintPlainText } from '../../../../lib/filter-dsl/printer';

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

  filter: Filter = null;
  error: string = null;

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
    console.log(this.editableContent.nativeElement.innerText.length, originalCaretPos);
    const originalText = this.editableContent.nativeElement.innerText;
    try {
      this.filter = Parse(this.editableContent.nativeElement.innerText);
      if (this.filter != null) {
        this.editableContent.nativeElement.innerHTML = '';
        this.renderer.appendChild(this.editableContent.nativeElement, PrintHTML(this.renderer, this.filter));
        this.clearError();
        this.moveCaretTo(originalCaretPos);
      }
    } catch (e) {
      console.error(e);
      this.filter = null;
      //this.editableContent.nativeElement.innerText = originalText;
      this.setError(e.message);
    }
  }

  clearError() {
    this.error = null;
  }

  setError(msg: string) {
    this.error = msg;
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
    const range = window.getSelection().getRangeAt(0);
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
