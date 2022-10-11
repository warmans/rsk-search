import { ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';
import { parse } from 'marked';

@Component({
  selector: 'app-markdown',
  templateUrl: './markdown.component.html',
  styleUrls: ['./markdown.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MarkdownComponent implements OnInit {

  @Input()
  set raw(value: string) {
    this._raw = value;
    this.render();
  }

  get raw(): string {
    return this._raw;
  }

  private _raw: string;

  renderedHTML: string;

  constructor() {
  }

  ngOnInit(): void {
  }

  render() {
    this.renderedHTML = parse(this.raw);
  }
}
