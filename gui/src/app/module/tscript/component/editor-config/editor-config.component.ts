import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { KeyPressEventCodes } from '../../../../lib/keys';

@Component({
  selector: 'app-editor-config',
  templateUrl: './editor-config.component.html',
  styleUrls: ['./editor-config.component.scss']
})
export class EditorConfigComponent implements OnInit {

  open: boolean;

  @Input()
  initialConfig: EditorConfig = new EditorConfig();

  @Output()
  configUpdated: EventEmitter<EditorConfig> = new EventEmitter<EditorConfig>();

  configForm: FormGroup = new FormGroup({
    'playPauseKey': new FormControl(),
    'backtrack': new FormControl(),
    'fastForwardKey': new FormControl(),
    'rewindKey': new FormControl(),
    'autoSeek': new FormControl(),
  });

  keyCodes = KeyPressEventCodes;

  constructor() {
  }

  ngOnInit(): void {
    this.configForm.get('playPauseKey').setValue(this.initialConfig.playPauseKey);
    this.configForm.get('backtrack').setValue(this.initialConfig.backtrack);
    this.configForm.get('fastForwardKey').setValue(this.initialConfig.fastForwardKey);
    this.configForm.get('rewindKey').setValue(this.initialConfig.rewindKey);
    this.configForm.get('autoSeek').setValue(this.initialConfig.autoSeek);
  }

  emitUpdatedconfig() {
    this.open = false;
    this.initialConfig = new EditorConfig(
      this.configForm.get('playPauseKey').value,
      this.configForm.get('backtrack').value,
      this.configForm.get('fastForwardKey').value,
      this.configForm.get('rewindKey').value,
      this.configForm.get('autoSeek').value,
    );
    this.configUpdated.next(this.initialConfig);
  }
}

export class EditorConfig {
  constructor(
    public playPauseKey: string = 'Insert',
    public backtrack: number = 3,
    public fastForwardKey: string = 'Pause',
    public rewindKey: string = 'ScrollLock',
    public autoSeek: boolean = true) {
  }
}
