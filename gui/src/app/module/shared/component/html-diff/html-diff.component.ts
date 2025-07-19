import { Component, Input, OnInit, ViewEncapsulation } from '@angular/core';
import { html } from 'diff2html';

@Component({
    selector: 'app-html-diff',
    templateUrl: './html-diff.component.html',
    encapsulation: ViewEncapsulation.None,
    styleUrls: [
        '../../../../../../node_modules/diff2html/bundles/css/diff2html.min.css',
        './html-diff.component.scss'
    ],
    standalone: false
})
export class HtmlDiffComponent implements OnInit {

  @Input()
  set unifiedDiff(value: string) {
    this._unifiedDiff = value;
    this.outputHtml = (value) ? html(value, { drawFileList: false, matching: 'lines' }) : '';
  }

  get unifiedDiff(): string {
    return this._unifiedDiff;
  }

  private _unifiedDiff: string;

  outputHtml: string;

  constructor() {
  }

  ngOnInit(): void {
  }

}
