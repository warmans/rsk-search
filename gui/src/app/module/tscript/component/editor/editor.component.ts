import { Component, ElementRef, EventEmitter, Input, OnInit, Output, Renderer2, ViewChild } from '@angular/core';
import { debounceTime, distinctUntilChanged } from 'rxjs/operators';

@Component({
  selector: 'app-editor',
  templateUrl: './editor.component.html',
  styleUrls: ['./editor.component.scss']
})
export class EditorComponent implements OnInit {

  _textContent: string = '';

  @Input()
  set textContent(f: string) {
    this._textContent = f;
    this.updateInnerHtml();
  }

  get textContent(): string {
    return this._textContent;
  }

  @ViewChild('textContainer')
  editableContent: ElementRef;

  @Output()
  contentUpdated: EventEmitter<string> = new EventEmitter<string>();

  constructor(private renderer: Renderer2) {
  }

  ngOnInit(): void {
  }

  updateInnerHtml() {
    if (!this._textContent) {
      return;
    }

    this.editableContent.nativeElement.innerText = "";

    const lines = this._textContent.split('\n');

    lines.forEach((line) => {
      if (line.match(/^#OFFSET:.*/g)){
        const el = this.renderer.createElement("span");
        el.innerText = `${line}\n`;
        el.className = 'do-not-edit';

        this.renderer.appendChild(this.editableContent.nativeElement, el);
      } else {
        const el = this.renderer.createElement("span");
        el.innerText = `${line}\n`;
        this.renderer.appendChild(this.editableContent.nativeElement, el);
      }
    });
  }

  onKeyup(event: KeyboardEvent) {
    this.contentUpdated.next(this.editableContent.nativeElement.innerText);
  }

}
