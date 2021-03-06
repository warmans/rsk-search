import { Component, EventEmitter, HostListener, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { Observable, Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged, takeUntil } from 'rxjs/operators';
import { SessionService } from '../../../core/service/session/session.service';
import { getFirstOffset, parseTranscript, Tscript } from '../../../shared/lib/tscript';
import { AudioPlayerComponent } from '../../../shared/component/audio-player/audio-player.component';
import { AlertService } from '../../../core/service/alert/alert.service';
import { EditorConfig, EditorConfigComponent } from '../../component/editor-config/editor-config.component';
import { formatDistance } from 'date-fns';
import { RskChunk, RskChunkContribution, RskContributionState } from '../../../../lib/api-client/models';

@Component({
  selector: 'app-submit',
  templateUrl: './submit.component.html',
  styleUrls: ['./submit.component.scss']
})
export class SubmitComponent implements OnInit, OnDestroy {

  authenticated: boolean = false;
  authError: string;
  chunk: RskChunk;
  contribution: RskChunkContribution;
  userCanEdit: boolean = true;
  userIsOwner: boolean = true;
  userIsApprover: boolean = false;

  // to stop the caret from getting messed up by updates we need to separate the input
  // data from the output.
  initialTranscript: string = '';
  updatedTranscript: string = '';
  firstOffset: number = -1;

  contentUpdated: Subject<string> = new Subject<string>();

  editorConfig: EditorConfig = localStorage.getItem('editor-config') ? JSON.parse(localStorage.getItem('editor-config')) as EditorConfig : new EditorConfig();

  audioPlayerURL: string = '';
  parsedTscript: Tscript;
  showHelp: boolean = false;

  cStates = RskContributionState;

  loading: boolean[] = [];
  lastUpdateTimestamp: Date;

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

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private apiClient: SearchAPIClient,
    private titleService: Title,
    private sessionService: SessionService,
    private alertService: AlertService,
  ) {
    titleService.setTitle('Contribute');

    route.paramMap.pipe(takeUntil(this.$destroy)).subscribe((d: Data) => {


      if (d.params['contribution_id']) {

        // load content from existing contribution
        this.loading.push(true);
        this.apiClient.getContribution({
          contributionId: d.params['contribution_id']
        }).pipe(takeUntil(this.$destroy)).subscribe((res: RskChunkContribution) => {
          this.setContribution(res);
        }).add(() => this.loading.shift());

        this.loading.push(true);
        this.apiClient.getChunk({ id: d.params['id'] }).pipe(takeUntil(this.$destroy)).subscribe(
          (v) => {
            if (!v) {
              return;
            }
            this.chunk = v;
            this.audioPlayerURL = `https://storage.googleapis.com/warmans-transcription-audio/${v.id}.mp3`;
          }
        ).add(() => this.loading.shift());

      } else {

        // load everything from the chunk

        this.loading.push(true);
        this.apiClient.getChunk({ id: d.params['id'] }).pipe(takeUntil(this.$destroy)).subscribe(
          (v) => {
            if (!v) {
              return;
            }
            titleService.setTitle(`Contribute :: ${v.id}`);

            this.chunk = v;
            this.audioPlayerURL = `https://storage.googleapis.com/warmans-transcription-audio/${v.id}.mp3`;

            this.setInitialTranscript(this.chunk.raw);
          }
        ).add(() => this.loading.shift());
      }

      sessionService.onTokenChange.pipe(takeUntil(this.$destroy)).subscribe((token: string): void => {
        if (token != null) {
          this.authenticated = true;
        }
      });
    });
  }

  ngOnDestroy(): void {
    this.$destroy.next();
    this.$destroy.complete();
  }

  ngOnInit(): void {
    this.contentUpdated.pipe(takeUntil(this.$destroy), distinctUntilChanged(), debounceTime(1000)).subscribe((v) => {
      this.updatedTranscript = v;
      this.backupContent(v);
      this.updatePreview(v);
      if (this.contribution && this.userCanEdit) {
        this.update();
      }
    });
  }

  setContribution(res: RskChunkContribution) {
    this.titleService.setTitle(`contribute :: ${res.chunkId} :: ${res.id}`);

    this.contribution = res;
    this.setInitialTranscript(res.transcript);

    this.userCanEdit = res.state === RskContributionState.STATE_PENDING;
    if (!this.sessionService.getClaims().approver) {
      this.userCanEdit = this.sessionService.getClaims()?.author_id === res.author.id;
    }

    this.userIsOwner = this.sessionService.getClaims()?.author_id === res.author.id || this.sessionService.getClaims().approver;
    this.userIsApprover = this.sessionService.getClaims().approver;
  }

  setInitialTranscript(text: string) {
    this.initialTranscript = this.getBackup() ? this.getBackup() : text;
    this.contentUpdated.next(this.initialTranscript);
    this.firstOffset = getFirstOffset(this.initialTranscript);
  }

  setUpdatedTranscript(text: string) {
    this.contentUpdated.next(text);
  }

  updatePreview(content: string) {
    this.parsedTscript = parseTranscript(content);
  }

  backupContent(text: string) {
    if (!this.chunk) {
      return;
    }
    localStorage.setItem(`chunk-backup-${(this.contribution) ? this.contribution.id : this.chunk.id}`, text);
  }

  getBackup(): string {
    return localStorage.getItem(`chunk-backup-${(this.contribution) ? this.contribution.id : this.chunk.id}`);
  }

  resetToRaw() {
    if (confirm('Really reset editor to raw raw transcript?')) {
      this.initialTranscript = this.contribution ? this.contribution.transcript : this.chunk.raw;
      this.updatePreview(this.initialTranscript);
    }
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
    return formatDistance(this.lastUpdateTimestamp, new Date());
  }

  create() {
    if (!this.contribution) {
      this.loading.push(true);
      this.apiClient.createChunkContribution({
        chunkId: this.chunk.id,
        body: { chunkId: this.chunk.id, transcript: this.updatedTranscript }
      }).pipe(takeUntil(this.$destroy)).subscribe((res: RskChunkContribution) => {
        this.backupContent(''); // clear backup so that the content always matches what was submitted.
        this.alertService.success('Created', 'Draft was created. It will now be auto-saved on change.');
        this.router.navigate(['/chunk', this.chunk.id, 'contrib', res.id]);
      }).add(() => this.loading.shift());
    }
  }

  update() {
    this._update(this.contribution.state).subscribe((res: RskChunkContribution) => {
      this.lastUpdateTimestamp = new Date();
    });
  }

  private _update(state: RskContributionState): Observable<RskChunkContribution> {
    return this.apiClient.updateContribution({
      contributionId: this.contribution.id,
      body: {
        contributionId: this.contribution.id,
        transcript: this.updatedTranscript,
        state: state
      }
    }).pipe(takeUntil(this.$destroy));
  }

  private _updateState(state: RskContributionState) {
    this.loading.push(true);
    this.apiClient.requestContributionState({
      contributionId: this.contribution.id,
      body: {
        contributionId: this.contribution.id,
        requestState: state,
      }
    }).pipe(takeUntil(this.$destroy)).subscribe((res) => {
      this.setContribution(res);
      switch (state) {
        case RskContributionState.STATE_PENDING:
          this.alertService.success('Retracted', 'Submission is now back in the pending state. It will not be reviewed until is is re-submitted.');
          return;
        case RskContributionState.STATE_APPROVED:
          this.alertService.success('Approved', 'Submission was approved.');
          return;
        case RskContributionState.STATE_REQUEST_APPROVAL:
          this.alertService.success('Submitted', 'Submission is now awaiting manual approval by an approver. This usually takes around 24 hours.');
          return;
        case RskContributionState.STATE_REJECTED:
          this.alertService.success('Rejected', 'Submission was rejected.');
          return;
      }
    }).add(() => this.loading.shift());
  }

  markComplete() {
    this.loading.push(true);
    this._update(RskContributionState.STATE_REQUEST_APPROVAL).subscribe((res: RskChunkContribution) => {
      this.setContribution(res);
      this.lastUpdateTimestamp = new Date();
      this.alertService.success('Submitted', 'Submission is now awaiting manual approval by an approver. This usually takes around 24 hours.');
    }).add(() => this.loading.shift());
  }

  markIncomplete() {
    this._updateState(RskContributionState.STATE_PENDING);
  }

  markApproved() {
    this._updateState(RskContributionState.STATE_APPROVED);
  }

  markRejected() {
    this._updateState(RskContributionState.STATE_REJECTED);
  }
}
