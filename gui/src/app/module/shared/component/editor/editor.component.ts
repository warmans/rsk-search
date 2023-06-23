import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  HostListener,
  Input,
  OnDestroy,
  OnInit,
  Output,
  ViewChild
} from '@angular/core';
import {EditorConfig, EditorConfigComponent} from '../editor-config/editor-config.component';
import {Subject} from 'rxjs';
import {getFirstOffset} from '../../lib/tscript';
import {distinctUntilChanged, takeUntil} from 'rxjs/operators';
import {formatDistance, fromUnixTime, getUnixTime, isBefore, subDays} from 'date-fns';
import {EditorInputComponent} from '../editor-input/editor-input.component';
import {AudioService, PlayerState, Status} from '../../../core/service/audio/audio.service';
import {FindReplace} from '../find-replace/find-replace.component';

const LOCAL_STORAGE_PREFIX = 'content-backup';

@Component({
  selector: 'app-editor',
  templateUrl: './editor.component.html',
  styleUrls: ['./editor.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class EditorComponent implements OnInit, OnDestroy {

  @Input()
  contentID: string;

  @Input()
  contentVersion: string;

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

  @ViewChild('editorConfigModal')
  editorConfigModal: EditorConfigComponent;

  @ViewChild('editorInput')
  editorComponent: EditorInputComponent;

  $destroy: EventEmitter<void> = new EventEmitter<void>();

  @HostListener('document:keydown', ['$event'])
  handleKeyboardEvent(event: KeyboardEvent): boolean {
    if (this._editorConfig.autoSeek === undefined ? true : this._editorConfig.autoSeek) {
      if (this._editorConfig?.playPauseKey && event.key === (this._editorConfig?.playPauseKey)) {
        this.audioService.toggleAudio(0 - (this._editorConfig?.backtrack || 3));
        return false;
      }
      if (this._editorConfig?.rewindKey && event.key === (this._editorConfig?.rewindKey)) {
        this.skipBackwards();
        return false;
      }
      if (this._editorConfig?.fastForwardKey && event.key === (this._editorConfig?.fastForwardKey)) {
        this.skipForward();
        return false;
      }
      if (this._editorConfig?.fastForwardKey && event.key === (this._editorConfig?.fastForwardKey)) {
        this.skipForward();
        return false;
      }
      if (this._editorConfig?.insertOffsetKey && event.key === (this._editorConfig?.insertOffsetKey)) {
        this.insertOffsetAboveCaret();
        return false;
      }
      if (this._editorConfig?.insertSynKey && event.key === (this._editorConfig?.insertSynKey)) {
        this.insertSynAboveCaret();
        return false;
      }
      return true;
    }
  }

  constructor(private audioService: AudioService, private cdr: ChangeDetectorRef) {
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

    try {
      this.cleanBackups();
    } catch (e) {
      console.error("failed to cleanup local storage");
    }

    this.contentUpdated.pipe(distinctUntilChanged(), takeUntil(this.$destroy)).subscribe((v) => {
      if (!v) {
        return;
      }
      try {
        this.backupContent(v);
      } catch (e) {
        console.error("cannot write to local storage, storage will be cleared.", e);
        localStorage.clear();
      }
      this.save();
    });
  }

  loadAudio(andPlay?: boolean) {
    const parts = (this._audioPlayerURL || '').split('/');
    const name = parts[parts.length - 1] ? parts[parts.length - 1] : null;
    if (name) {
      this.audioService.setAudioSrc(name, null, this._audioPlayerURL, true);
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
    if (this.editorComponent) {
      this.editorComponent.textContent = this.initialTranscript;
    }
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
    localStorage.setItem(this.localBackupKey(), text);
  }

  getBackup(): string {
    if (!this.contentID || !this._allowEdit) {
      return;
    }
    return localStorage.getItem(this.localBackupKey());
  }

  clearBackup(): void {
    localStorage.removeItem(this.localBackupKey());
    this.fromBackup = false;
  }

  localBackupKey(): string {
    return `${LOCAL_STORAGE_PREFIX}-${getUnixTime(new Date())}-${this.contentID}${this.contentVersion ? '-' + this.contentVersion : ''}`;
  }

  // local storage cannot be allowed to fill up.
  // Remove backups if they're either an old format, or older than 1 week.
  cleanBackups() {
    for (let i = 0; i < localStorage.length; i++) {
      const key: string = localStorage.key(i);
      if (key.startsWith(LOCAL_STORAGE_PREFIX)) {
        const parts: string[] = key.replace(LOCAL_STORAGE_PREFIX + "-", "").split("-");
        if (parts.length > 0 && (/^[0-9]+/).test(parts[0])) {
          const itemDate = fromUnixTime(parseInt(parts[0]));
          if (isBefore(itemDate, subDays(new Date(), 7))) {
            // remove backups older than a week
            localStorage.removeItem(key);
          }
        } else {
          // remove legacy backups
          localStorage.removeItem(key);
        }
      }
    }
  }

  resetToRaw() {
    if (confirm('Really reset editor to raw transcript?')) {
      this.clearBackup();
      this.setInitialTranscript(this._rawTranscript);
      this.cdr.detectChanges();
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
    this.editorComponent.insertOffsetAboveCaret(Math.floor(startOffset + (this.audioStatus?.currentTime || 0) - (this._editorConfig.insertOffsetBacktrack || 0)));
  }

  insertSynAboveCaret() {
    this.insertTextAboveCaret('#SYN: ');
  }

  insertTextAboveCaret(text: string) {
    this.editorComponent.insertTextAboveCaret(text);
  }

  refreshEditorHTML() {
    this.editorComponent.refreshInnerHtml();
  }

  runFindAndReplace(vals: FindReplace) {
    this.editorComponent.findAndReplace(vals.find, vals.replace);
  }
}
