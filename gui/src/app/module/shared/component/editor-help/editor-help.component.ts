import { ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';
import { EditorConfig } from '../editor-config/editor-config.component';

@Component({
    selector: 'app-editor-help',
    templateUrl: './editor-help.component.html',
    styleUrls: ['./editor-help.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class EditorHelpComponent implements OnInit {

  @Input()
  editorConfig: EditorConfig;

  constructor() {
  }

  ngOnInit(): void {
  }

}
