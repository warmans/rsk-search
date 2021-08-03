import { Component, EventEmitter, HostListener, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { EditorConfig, EditorConfigComponent } from '../editor-config/editor-config.component';
import { Subject } from 'rxjs';
import { getFirstOffset } from '../../lib/tscript';
import { AudioPlayerComponent } from '../audio-player/audio-player.component';
import { debounceTime, distinctUntilChanged, takeUntil } from 'rxjs/operators';
import { formatDistance } from 'date-fns';

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
  audioPlayerURL: string = '';

  @Input()
  allowEdit: boolean = false;

  @Input()
  isSaved: boolean = false;

  @Output()
  handleSave: EventEmitter<string> = new EventEmitter<string>();

  fromBackup: boolean = false;

  editorConfig: EditorConfig = localStorage.getItem('editor-config') ? JSON.parse(localStorage.getItem('editor-config')) as EditorConfig : new EditorConfig();

  contentUpdated: Subject<string> = new Subject<string>();

  initialTranscript: string = '';
  updatedTranscript: string = '';
  firstOffset: number = -1;

  showHelp: boolean = false;

  $destroy: EventEmitter<boolean> = new EventEmitter<boolean>();

  @ViewChild('audioPlayer')
  audioPlayer: AudioPlayerComponent;

  @ViewChild('editorConfigModal')
  editorConfigModal: EditorConfigComponent;

  @HostListener('document:keydown', ['$event'])
  handleKeyboardEvent(event: KeyboardEvent): boolean {
    if (this.editorConfig.autoSeek === undefined ? true : this.editorConfig.autoSeek) {
      if (event.key === (this.editorConfig?.playPauseKey || 'Insert')) {
        this.audioPlayer.toggle(0 - (this.editorConfig?.backtrack || 3));
        return false;
      }
      if (event.key === (this.editorConfig?.rewindKey || 'ScrollLock')) {
        this.audioPlayer.play(-3);
        return false;
      }
      if (event.key === (this.editorConfig?.fastForwardKey || 'Pause')) {
        this.audioPlayer.play(3);
        return false;
      }
      return true;
    }
  }

  constructor() {
  }

  ngOnDestroy(): void {
    this.$destroy.next();
    this.$destroy.complete();
  }

  ngOnInit(): void {
    this.contentUpdated.pipe(takeUntil(this.$destroy), distinctUntilChanged(), debounceTime(1000)).subscribe((v) => {
      this.updatedTranscript = v;
      this.backupContent(v);
      console.log("updated raw")
      //this.updatePreview(v);
      this.save();
    });
  }

  setInitialTranscript(text: string) {
    const backup = this.getBackup()
    this.initialTranscript = backup ? backup : text;
    if (backup) {
      this.fromBackup = true;
      console.log(backup);
    }
    this.contentUpdated.next(this.initialTranscript);
    this.firstOffset = getFirstOffset(this.initialTranscript);
  }

  backupContent(text: string) {
    if (!this.updatedTranscript || !this.contentID) {
      return;
    }
    localStorage.setItem(`content-backup-${this.contentID}`, text);
  }

  getBackup(): string {
    if (!this.contentID) {
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
      //this.updatePreview(this.initialTranscript);
    }
  }

  setUpdatedTranscript(text: string) {
    this.contentUpdated.next(text);
  }

  handleOffsetNavigate(offset: number) {
    if (this.editorConfig?.autoSeek === undefined ? true : this.editorConfig.autoSeek) {
      if (offset - this.firstOffset >= 0) {
        this.audioPlayer.seek(offset - this.firstOffset);
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
    this.editorConfig = cfg;
    localStorage.setItem('editor-config', JSON.stringify(cfg));
  }

  timeSinceSave(): string {
    return formatDistance(this.lastUpdateDate, new Date());
  }

  save(): void {
    if (this.handleSave) {
      this.handleSave.next(this.updatedTranscript);
    }
  }
}
