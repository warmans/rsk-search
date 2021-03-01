import { Component, ElementRef, EventEmitter, OnInit, Output, Renderer2, ViewChild } from '@angular/core';
import { Scan, Tok } from '../../../../lib/filter-dsl/scanner';

@Component({
  selector: 'app-search-bar',
  templateUrl: './search-bar.component.html',
  styleUrls: ['./search-bar.component.scss']
})
export class SearchBarComponent implements OnInit {

  @Output()
  executeQuery: EventEmitter<string> = new EventEmitter<string>();

  @ViewChild('editableContent')
  editableContent: ElementRef;

  keyboardEvents: EventEmitter<KeyboardEvent> = new EventEmitter();

  invalid = false;
  inputActive = false;
  simpleSearch: boolean = false;

  // current state of element
  tokens: Tok[] = [];
  activeToken: Tok = null;

  // // predicts what the user must input next
  predictedNextToken: Tok = null;

  caretPos: number;
  contentLength: number;

  constructor(private renderer: Renderer2) {
  }

  ngOnInit(): void {
  }

  emit() {
    if (this.simpleSearch) {
      const q = this.editableContent.nativeElement.innerText.replace(/["]+/g, '');

      //todo: I think it should check for quotes. If any text is quoted extract it and add it as another condition
      // e.g. 'foo bar "baz"' becomes content ~= "foo bar" and content = "baz"
      this.executeQuery.emit(`content ~= "${q}"`);
    } else {
      this.executeQuery.emit(this.tokensAsText());
    }
  }

  toggleAdvanced(ev: any) {
    this.simpleSearch = !ev.target.checked;
    if (this.simpleSearch) {
      this.invalid = false;
      this.editableContent.nativeElement.innerHTML = this.tokensAsText();
    } else {
      this.update();
    }
  }

  onFocus() {
    this.inputActive = true;
  }

  onClick() {
    this.inputActive = true;
    this.update();
  }

  onKeydown(key: KeyboardEvent): boolean {
    switch (key.code) {
      case 'ArrowDown':
      case 'ArrowUp':
      case 'Enter':
        this.inputActive = true;
        key.preventDefault();
        this.emit();
        return false;
      case 'Escape':
        this.inputActive = false;
    }
  }

  onKeypress(evt: KeyboardEvent) {

    // todo: should only update if the key is a valid one (i.e. not shift or something)
    this.update();

    // forward event to any other components e.g. value-list
    this.keyboardEvents.next(evt);
  }

  update() {
    if (this.simpleSearch) {
      return;
    }

    if (!this.parseInput()) {
      return;
    }

    this.saveCaretPosition(this.editableContent.nativeElement);
    this.findActiveElements();

    // // update UI
    this.render();
    this.moveCaretTo(this.caretPos);
  }

  parseInput(): boolean {
    this.invalid = false;
    try {
      this.tokens = Scan(this.editableContent.nativeElement.textContent);
    } catch {
      this.invalid = true;
    }
    return !this.invalid;
  }

  findActiveElements() {
    this.activeToken = null;

    this.tokens.forEach(tok => {
      if (this.caretPos - 1 >= tok.start && this.caretPos - 1 <= tok.end) {
        this.activeToken = tok;
      }
    });
  }

  render() {
    const rendered = this.renderer.createElement('span');

    this.tokens.forEach((tok: Tok) => {
      const tokEl = this.renderer.createElement('span');
      tokEl.className = tok.tag;
      tokEl.textContent = tok.lexeme;
      this.renderer.appendChild(rendered, tokEl);
    });

    // clear input
    this.editableContent.nativeElement.innerHTML = '';

    // update with styled content
    this.renderer.appendChild(this.editableContent.nativeElement, rendered);
  }

  moveCaretTo(position: number) {

    // move to end (no idea why it doesn't work in the normal way)
    if (this.editableContent.nativeElement.textContent.length === position) {
      let range, selection;
      range = document.createRange();
      range.selectNodeContents(this.editableContent.nativeElement);
      range.collapse(false);
      selection = window.getSelection();
      selection.removeAllRanges();
      selection.addRange(range);
      return;
    }

    const node = getTextNodeAtPosition(
      this.editableContent.nativeElement,
      position
    );
    const sel = window.getSelection();
    sel.collapse(node.node, node.position);
  }

  saveCaretPosition(context) {
    const range = window.getSelection().getRangeAt(0);
    const selected = range.toString().length; // *
    const preCaretRange = range.cloneRange();
    preCaretRange.selectNodeContents(context);
    preCaretRange.setEnd(range.endContainer, range.endOffset);

    if (selected) {
      this.caretPos = preCaretRange.toString().length - selected;
    } else {
      this.caretPos = preCaretRange.toString().length;
    }
    this.contentLength = context.textContent.length;
  }

  tokensAsText(): string {
    const t: string[] = [];
    this.tokens.forEach((v) => {
      t.push(v.lexeme);
    });
    return t.join('');
  }

  // updateActiveTokenValue(value: string[]) {
  //   if (!this.activeToken) {
  //     return;
  //   }
  //   const strVal = value.join(',');
  //   this.editableContent.nativeElement.textContent = spliceString(
  //     this.editableContent.nativeElement.textContent,
  //     this.activeToken.start,
  //     this.activeToken.end,
  //     strVal,
  //   );
  //
  //   // move caret to end of the new token
  //   this.moveCaretTo(this.activeToken.start + strVal.length);
  //
  //   this.update();
  // }

  // appendToken(value: string[]) {
  //   this.editableContent.nativeElement.textContent += value.join(',');
  //
  //   // move caret to end of the new token
  //   this.moveCaretTo(this.editableContent.nativeElement.textContent.length);
  //   this.update();
  // }
}

function getTextNodeAtPosition(root, index) {
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

function spliceString(str: string, start: number, end: number, replace: string) {
  if (start < 0) {
    start = str.length + start;
    start = start < 0 ? 0 : start;
  }
  return str.slice(0, start) + (replace || '') + str.slice(end);
}
