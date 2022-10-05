import { AfterViewInit, Component, ElementRef, EventEmitter, Input, OnDestroy, OnInit, Output, Renderer2, ViewChild } from '@angular/core';
import { getOffsetValueFromLine, isOffsetLine } from '../../lib/tscript';
import { Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged } from 'rxjs/operators';

@Component({
  selector: 'app-editor',
  templateUrl: './editor.component.html',
  styleUrls: ['./editor.component.scss']
})
export class EditorComponent implements OnInit, OnDestroy, AfterViewInit {

  _textContent = '';

  @Input()
  set textContent(f: string) {
    this._textContent = f;
    if (this.editableContent) {
      this.updateInnerHtml(this._textContent);
    }
  }

  get textContent(): string {
    // do not return the innerText directly. If the div is hidden by the parent
    // before calling this, it will get text with no line breaks.
    return this._textContent;
  }

  @Input()
  readonly: boolean = false;

  @Input()
  wrap: boolean = false;

  @Output()
  textContentChange: EventEmitter<string> = new EventEmitter<string>();

  textChangeDebouncer: Subject<string> = new Subject<string>();

  @Output()
  atOffsetMarker: EventEmitter<number> = new EventEmitter<number>();

  @ViewChild('textContainer')
  editableContent: ElementRef;

  destory$: EventEmitter<any> = new EventEmitter<any>();

  constructor(private renderer: Renderer2) {
  }

  ngOnInit(): void {
    this.textChangeDebouncer.pipe(distinctUntilChanged(), debounceTime(500)).subscribe((v: string) => {
      this.textContentChange.next(v);
    });
  }

  onKeypress() {
    this.contentChanged();
    this.tryEmitOffset();
    //this.tryShowAutocomplete();
  }

  contentChanged() {
    this._textContent = this.editableContent.nativeElement.innerText;
    this.textChangeDebouncer.next(this._textContent);
  }

  ngOnDestroy(): void {
    this.destory$.next(true);
    this.destory$.complete();
  }

  refreshInnerHtml() {
    this.updateInnerHtml(this._textContent);
  }

  findAndReplace(find: string, replace: string) {

    // this is really dumb. If you want to preserve the undo buffer then you need top use this API... which is deprecated
    // and there is no replacement. And it sucks. Better than not allowing undo... I guess.
    // Currently, replace could break the HTML. I guess I should edit the innerText then re-create the inner HTML.
    // can't be bothered.
    this.editableContent.nativeElement.focus();
    document.execCommand('selectAll');
    document.execCommand(
      'insertHTML',
      false,
      this.editableContent.nativeElement.innerHTML.replace(new RegExp(`${find}`, 'g'), replace),
    );

    this.contentChanged();
  }

  updateInnerHtml(value: string) {
    if (!value || !this.editableContent) {
      return;
    }
    this.editableContent.nativeElement.innerText = '';

    const lines = value.split('\n');

    lines.forEach((line: string) => {
      line = line.trim();
      if (line.length === 0) {
        return;
      }
      if (line.match(/^#OFFSET:.*/g)) {
        this.renderer.appendChild(this.editableContent.nativeElement, this.newOffsetElement(line));
      } else if (line.match(/^#[\/]?(SYN|TRIVIA).*/g)) {
        this.renderer.appendChild(this.editableContent.nativeElement, this.newTextElement(line));
      } else {
        const el = this.renderer.createElement('span');
        el.innerText = `${line}\n`;
        this.renderer.appendChild(this.editableContent.nativeElement, el);
      }
    });
  }

  private newOffsetElement(line: string): any {
    const el = this.renderer.createElement('span');
    el.innerText = `${line}\n`;
    el.className = 'do-not-edit';
    return el;
  }

  private newTextElement(text: string): any {
    const el = this.renderer.createElement('span');
    el.innerText = `${text}\n`;
    el.className = 'meta-tag';
    return el;
  }

  private newAutocomplete(items: string[]) {
    const el = this.renderer.createElement('span');
    el.innerText = `${items.join(',')}\n`;
    el.className = 'autocomplete';
    return el;
  }

  private tryEmitOffset() {
    const caretFocus = this.getCaretFocus();
    if (caretFocus && isOffsetLine(caretFocus.line)) {
      this.atOffsetMarker.next(getOffsetValueFromLine(caretFocus.line));
    }
  }

  // private activeAutocomplete: {el: HTMLElement, parent: any};
  //
  // private tryShowAutocomplete() {
  //   const caretFocus = this.getCaretFocus();
  //   if (!caretFocus) {
  //     return;
  //   }
  //   if (this.activeAutocomplete) {
  //     this.renderer.removeChild(this.activeAutocomplete.parent, this.activeAutocomplete.el);
  //     this.activeAutocomplete = undefined;
  //   }
  //   if (!lineHasActorPrefix(caretFocus.line)) {
  //     let sel = document.getSelection();
  //     let nd = sel.anchorNode;
  //     let newNode = this.newAutocomplete(["foo", "bar"]);
  //     this.renderer.appendChild(nd, newNode);
  //     this.activeAutocomplete = {el: newNode, parent: nd};
  //   }
  // }

  insertOffsetAboveCaret(seconds: number): void {
    let sel = document.getSelection();
    let nd = sel.anchorNode;

    if (!this.nodeIsChildOfEditor(nd.parentElement)) {
      return;
    }
    this.renderer.insertBefore(nd.parentNode, this.newOffsetElement(`#OFFSET: ${seconds}`), nd);
    this.contentChanged();
  }

  insertTextAboveCaret(text: string): void {
    let sel = document.getSelection();
    let nd = sel.anchorNode;
    if (!this.nodeIsChildOfEditor(nd.parentElement)) {
      return;
    }
    this.renderer.insertBefore(nd.parentNode, this.newTextElement(text), nd);
    this.contentChanged();
  }

  private nodeIsChildOfEditor(el: HTMLElement): boolean {
    if (el === this.editableContent.nativeElement) {
      return true;
    }
    if (el.parentElement !== null) {
      return this.nodeIsChildOfEditor(el.parentElement);
    }
    return false;
  }

  getCaretFocus(): CaretFocus {
    let sel = document.getSelection();
    let nd = sel.anchorNode;
    if (!this.nodeIsChildOfEditor(nd.parentElement)) {
      return;
    }
    return new CaretFocus(nd.textContent, sel.focusOffset);
  }

  ngAfterViewInit(): void {
    this.updateInnerHtml(this._textContent);
  }
}

class CaretFocus {
  constructor(public line: string, public offset: number) {
  }
}
