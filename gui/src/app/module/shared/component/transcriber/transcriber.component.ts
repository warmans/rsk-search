import { Component, EventEmitter, HostListener, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { EditorConfig, EditorConfigComponent } from '../editor-config/editor-config.component';
import { Subject } from 'rxjs';
import { getFirstOffset, Tscript } from '../../lib/tscript';
import { distinctUntilChanged, takeUntil } from 'rxjs/operators';
import { formatDistance } from 'date-fns';
import { EditorComponent } from '../editor/editor.component';
import { AudioService, PlayerState, Status } from '../../../core/service/audio/audio.service';

@Component({
  selector: 'app-transcriber',
  templateUrl: './transcriber.component.html',
  styleUrls: ['./transcriber.component.scss']
})
export class TranscriberComponent implements OnInit, OnDestroy {

  @Input()
  contentID: string;

  @Input()
  set rawTranscript(value: string) {
    this._rawTranscript = value;
    this.setInitialTranscript(value);
  }

  get rawTranscript(): string {
    return this._rawTranscript;
  }

  private _rawTranscript: string = '';

  @Input()
  lastUpdateDate: Date;

  @Input()
  set audioPlayerURL(value: string) {
    this._audioPlayerURL = value;
    this.loadAudio();
  }

  get audioPlayerURL(): string {
    return this._audioPlayerURL;
  }

  private _audioPlayerURL: string = '';

  @Input()
  set allowEdit(value: boolean) {
    this._allowEdit = value;
    if (!value) {
      this.activeTab = 'preview';
    }
  }

  get allowEdit(): boolean {
    return this._allowEdit;
  }

  private _allowEdit: boolean = false;

  @Input()
  isSaved: boolean = false;

  @Output()
  handleSave: EventEmitter<string> = new EventEmitter<string>();

  @Output()
  activateTab: EventEmitter<string> = new EventEmitter<string>();

  @Input()
  enableDiff: boolean;

  @Input()
  unifiedDiff: string;

  // when editing a chunk we need to get the first offset and use it to modify the current offset
  // since we are working with a random chunk of audio from the episode.
  @Input()
  chunkMode: boolean = false;

  fromBackup: boolean = false;

  get editorConfig(): EditorConfig {
    return this._editorConfig;
  }

  set editorConfig(cfg: EditorConfig) {
    this._editorConfig = cfg;
    this.audioService.setPlaybackRate(cfg.playbackRate || 1);
    localStorage.setItem('editor-config', JSON.stringify(cfg));
  }
  private _editorConfig: EditorConfig;

  contentUpdated: Subject<string> = new Subject<string>();

  initialTranscript: string = '';
  firstOffset: number = -1;

  showHelp: boolean = false;

  audioStatus: Status;
  playerStates = PlayerState;

  set activeTab(value: 'edit' | 'preview' | 'diff') {
    this.activateTab.next(value);
    this._activeTab = value;
  }

  get activeTab(): 'edit' | 'preview' | 'diff' {
    return this._activeTab;
  }

  private _activeTab: 'edit' | 'preview' | 'diff' = 'edit';

  parsedTscript: Tscript;

  @ViewChild('editorConfigModal')
  editorConfigModal: EditorConfigComponent;

  @ViewChild('editor')
  editorComponent: EditorComponent;

  $destroy: EventEmitter<boolean> = new EventEmitter<boolean>();

  @HostListener('document:keydown', ['$event'])
  handleKeyboardEvent(event: KeyboardEvent): boolean {
    if (this._editorConfig.autoSeek === undefined ? true : this._editorConfig.autoSeek) {
      if (event.key === (this._editorConfig?.playPauseKey || 'Insert')) {
        this.audioService.toggleAudio(0 - (this._editorConfig?.backtrack || 3));
        return false;
      }
      if (event.key === (this._editorConfig?.rewindKey || 'ScrollLock')) {
        this.skipBackwards();
        return false;
      }
      if (event.key === (this._editorConfig?.fastForwardKey || 'Pause')) {
        this.skipForward()
        return false;
      }
      if (event.key === (this._editorConfig?.fastForwardKey || 'Pause')) {
        this.skipForward()
        return false;
      }
      if (event.key === (this._editorConfig?.insertOffsetKey || 'PrintScreen')) {
        this.insertOffsetAboveCaret();
        return false;
      }
      return true;
    }
  }

  constructor(public audioService: AudioService) {
    audioService.status.pipe(takeUntil(this.$destroy)).subscribe((sta) => {
      this.audioStatus = sta;
    });

    this.editorConfig = localStorage.getItem('editor-config') ? JSON.parse(localStorage.getItem('editor-config')) as EditorConfig : new EditorConfig();
  }

  ngOnDestroy(): void {
    this.audioService.reset();
    this.$destroy.next();
    this.$destroy.complete();
  }

  ngOnInit(): void {
    this.contentUpdated.pipe(takeUntil(this.$destroy), distinctUntilChanged()).subscribe((v) => {
      if (!v) {
        return;
      }
      this.backupContent(v);
      this.save();
    });
  }

  loadAudio(andPlay?: boolean) {
    const parts = (this._audioPlayerURL || '').split('/');
    const name = parts[parts.length - 1] ? parts[parts.length - 1] : null;
    if (name) {
      this.audioService.setAudioSrc(name, this._audioPlayerURL, true);
      if (andPlay) {
        this.audioService.playAudio();
      }
    }
  }

  skipForward() {
    this.audioService.playAudio(3);
  }

  skipBackwards() {
    this.audioService.playAudio(-3);
  }

  togglePlayer() {
    this.audioService.toggleAudio();
  }

  getContentSnapshot(): string {
    return `${this.editorComponent?.textContent || ''}`;
  }

  setInitialTranscript(text: string) {
    const backup = this.getBackup();
    this.initialTranscript = backup ? backup : text;
    if (backup) {
      this.fromBackup = true;
    }
    this.contentUpdated.next(this.initialTranscript);
    this.firstOffset = this.chunkMode ? getFirstOffset(this.initialTranscript) : 0;
  }

  backupContent(text: string) {
    if (!text || !this.contentID || !this._allowEdit) {
      return;
    }
    localStorage.setItem(`content-backup-${this.contentID}`, text);
  }

  getBackup(): string {
    if (!this.contentID || !this._allowEdit) {
      return;
    }
    return localStorage.getItem(`content-backup-${this.contentID}`);
  }

  clearBackup(): void {
    localStorage.removeItem(`content-backup-${this.contentID}`);
    this.fromBackup = false;
  }

  resetToRaw() {
    if (confirm('Really reset editor to raw raw transcript?')) {
      this.clearBackup();
      this.setInitialTranscript(this._rawTranscript);
    }
  }

  handleContentUpdated() {
    this.contentUpdated.next(this.getContentSnapshot());
  }

  handleOffsetNavigate(offset: number) {
    if (this._editorConfig?.autoSeek === undefined ? true : this._editorConfig.autoSeek) {
      if (offset - this.firstOffset >= 0) {
        this.audioService.seekAudio(offset - this.firstOffset);
      }
    }
  }

  openEditorConfig() {
    if (!this.editorConfigModal) {
      return;
    }
    this.editorConfigModal.open = true;
  }

  handleEditorConfigUpdated(cfg: EditorConfig) {

  }

  timeSinceSave(): string {
    return formatDistance(this.lastUpdateDate, new Date());
  }

  save(): void {
    if (this.handleSave) {
      this.handleSave.next(this.getContentSnapshot());
    }
  }

  insertOffsetAboveCaret() {
    let startOffset = this.firstOffset > -1 ? this.firstOffset : 0;
    this.editorComponent.insertOffsetAboveCaret(Math.floor(startOffset + (this.audioStatus?.currentTime || 0)-0.5));
  }

  insertTextAboveCaret(text: string) {
    this.editorComponent.insertTextAboveCaret(text);
  }


}
