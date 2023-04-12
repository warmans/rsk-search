import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { UntypedFormControl, UntypedFormGroup, ValidationErrors, ValidatorFn } from '@angular/forms';
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

  configForm: UntypedFormGroup = new UntypedFormGroup({
    'playPauseKey': new UntypedFormControl(),
    'playbackRate': new UntypedFormControl(),
    'backtrack': new UntypedFormControl(),
    'fastForwardKey': new UntypedFormControl(),
    'rewindKey': new UntypedFormControl(),
    'insertOffsetKey': new UntypedFormControl(),
    'insertOffsetBacktrack': new UntypedFormControl(),
    'insertSynKey': new UntypedFormControl(),
    'autoSeek': new UntypedFormControl(),
    'wrapText': new UntypedFormControl(),
  });

  keyCodes = KeyPressEventCodes;

  constructor() {
    this.configForm.setValidators(this.validateKeyBindings());
  }

  ngOnInit(): void {
    this.configForm.get('playPauseKey').setValue(this.initialConfig.playPauseKey);
    this.configForm.get('playbackRate').setValue(this.initialConfig.playbackRate || 1.0);
    this.configForm.get('backtrack').setValue(this.initialConfig.backtrack);
    this.configForm.get('fastForwardKey').setValue(this.initialConfig.fastForwardKey);
    this.configForm.get('rewindKey').setValue(this.initialConfig.rewindKey);
    this.configForm.get('insertOffsetKey').setValue(this.initialConfig.insertOffsetKey);
    this.configForm.get('insertOffsetBacktrack').setValue(this.initialConfig.insertOffsetBacktrack || 0);
    this.configForm.get('insertSynKey').setValue(this.initialConfig.insertSynKey);
    this.configForm.get('autoSeek').setValue(this.initialConfig.autoSeek);
    this.configForm.get('wrapText').setValue(this.initialConfig.wrapText === undefined ? true : this.initialConfig.wrapText);
  }

  validateKeyBindings(): ValidatorFn {
    return (group: UntypedFormGroup): ValidationErrors => {
      const boundKeys: { [index: string]: boolean } = {};

      for (let ctrlName in group.controls) {
        if (!group.controls.hasOwnProperty(ctrlName) || group.get(ctrlName).value === '') {
          continue;
        }
        if (KeyPressEventCodes.indexOf(group.get(ctrlName).value) > -1 && boundKeys[group.get(ctrlName).value]) {
          group.get(ctrlName).setErrors({ 'duplicateValue': [`${group.get(ctrlName).value} already in use`] });
        }
        boundKeys[group.get(ctrlName).value] = true;
      }
      return;
    };
  }

  emitUpdatedconfig() {
    this.open = false;
    this.initialConfig = new EditorConfig(
      this.configForm.get('playPauseKey').value,
      this.configForm.get('playbackRate').value,
      this.configForm.get('backtrack').value,
      this.configForm.get('fastForwardKey').value,
      this.configForm.get('rewindKey').value,
      this.configForm.get('insertOffsetKey').value,
      this.configForm.get('insertOffsetBacktrack').value,
      this.configForm.get('insertSynKey').value,
      this.configForm.get('autoSeek').value,
      this.configForm.get('wrapText').value,
    );
    this.configUpdated.next(this.initialConfig);
  }
}

export class EditorConfig {
  constructor(
    public playPauseKey: string = 'Insert',
    public playbackRate: number = 1.0,
    public backtrack: number = 3,
    public fastForwardKey: string = 'Pause',
    public rewindKey: string = 'ScrollLock',
    public insertOffsetKey: string = 'PrintScreen',
    public insertOffsetBacktrack: number = 0.5,
    public insertSynKey: string = '',
    public autoSeek: boolean = true,
    public wrapText: boolean = true) {
  }
}
