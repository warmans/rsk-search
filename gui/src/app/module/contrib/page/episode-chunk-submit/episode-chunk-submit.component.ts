import { Component, EventEmitter, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { Observable, Subject } from 'rxjs';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { SessionService } from '../../../core/service/session/session.service';
import { getFirstOffset } from '../../../shared/lib/tscript';
import { AlertService } from '../../../core/service/alert/alert.service';
import { RskChunk, RskChunkContribution, RskContributionState } from 'src/app/lib/api-client/models';
import { EditorComponent } from '../../../shared/component/editor/editor.component';

@Component({
    selector: 'app-episode-chunk-submit',
    templateUrl: './episode-chunk-submit.component.html',
    styleUrls: ['./episode-chunk-submit.component.scss'],
    standalone: false
})
export class EpisodeChunkSubmit implements OnInit, OnDestroy {

  authenticated: boolean = false;
  chunk: RskChunk;
  contribution: RskChunkContribution;
  userCanEdit: boolean = true;
  userIsOwner: boolean = true;
  userIsApprover: boolean = false;

  // to stop the caret from getting messed up by updates we need to separate the input
  // data from the output.
  initialTranscript: string = '';

  firstOffset: number = -1;

  contentUpdated: Subject<string> = new Subject<string>();

  audioPlayerURL: string = '';

  cStates = RskContributionState;

  loading: boolean[] = [];

  lastUpdateTimestamp: Date;

  activeTab: 'edit' | 'preview' | 'diff' = 'edit';

  rejectCallback: (contributionId: string, comment: string) => void = (contributionId: string, comment: string) => {
    this.markRejected(comment);
  };

  $destroy: EventEmitter<void> = new EventEmitter<void>();

  @ViewChild('editor')
  editor: EditorComponent;

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
        this.apiClient.getChunkContribution({
          contributionId: d.params['contribution_id']
        }).pipe(takeUntil(this.$destroy)).subscribe((res: RskChunkContribution) => {
          this.setContribution(res);
        }).add(() => this.loading.shift());

        this.loading.push(true);
        this.apiClient.getTranscriptChunk({ id: d.params['id'] }).pipe(takeUntil(this.$destroy)).subscribe(
          (v) => {
            if (!v) {
              return;
            }
            this.chunk = v;
          }
        ).add(() => this.loading.shift());

      } else {

        // load everything from the chunk

        this.loading.push(true);
        this.apiClient.getTranscriptChunk({ id: d.params['id'] }).pipe(takeUntil(this.$destroy)).subscribe(
          (v) => {
            if (!v) {
              return;
            }
            titleService.setTitle(`Contribute :: ${v.id}`);

            this.chunk = v;

            this.setInitialTranscript(this.chunk.raw);
          }
        ).add(() => this.loading.shift());
      }
    });

    sessionService.onTokenChange.pipe(takeUntil(this.$destroy)).subscribe((token: string): void => {
      if (token != null) {
        this.authenticated = true;
      }
    });
  }

  ngOnDestroy(): void {
    this.$destroy.next();
    this.$destroy.complete();
  }

  ngOnInit(): void {
    this.contentUpdated.pipe(debounceTime(1000), takeUntil(this.$destroy)).subscribe((v) => {
      this.executeUpdate(v);
    });
  }

  executeUpdate(text: string): void {
    if (this.contribution && this.userCanEdit) {
      this.update();
    }
  }

  setContribution(res: RskChunkContribution) {
    this.titleService.setTitle(`Contribute :: ${res.chunkId} :: ${res.id}`);

    this.contribution = res;

    // each time the contribution is saved the backup needs to be cleared to ensure the editor is in sync with the
    // server-side version of the transcript.
    this.editor.clearBackup();

    this.setInitialTranscript(res.transcript);

    this.userCanEdit = res.state === RskContributionState.STATE_PENDING;
    if (!this.sessionService.getClaims().approver) {
      this.userCanEdit = this.sessionService.getClaims()?.author_id === res.author.id;
    }

    this.userIsOwner = this.sessionService.getClaims()?.author_id === res.author.id || this.sessionService.getClaims().approver;
    this.userIsApprover = this.sessionService.getClaims().approver;
  }

  setInitialTranscript(text: string) {
    this.initialTranscript = text;
    this.contentUpdated.next(this.initialTranscript);
    this.firstOffset = getFirstOffset(this.initialTranscript);
  }

  handleSave(text: string) {
    this.contentUpdated.next(text);
  }

  create() {
    if (!this.contribution) {
      this.loading.push(true);
      this.apiClient.createChunkContribution({
        chunkId: this.chunk.id,
        body: {
          transcript: this.editor.getContentSnapshot()
        }
      }).pipe(takeUntil(this.$destroy)).subscribe((res: RskChunkContribution) => {
        this.editor.clearBackup();
        this.alertService.success('Created', 'Draft was created. It will now be auto-saved on change.');
        this.router.navigate(['/chunk', this.chunk.id, 'contrib', res.id]);
      }).add(() => this.loading.shift());
    }
  }

  update() {
    this._update(this.contribution.state).subscribe((res: RskChunkContribution) => {
      this.lastUpdateTimestamp = new Date();
      this.editor.clearBackup();
    });
  }

  private _update(state: RskContributionState): Observable<RskChunkContribution> {
    return this.apiClient.updateChunkContribution({
      contributionId: this.contribution.id,
      body: {
        transcript: this.editor.getContentSnapshot(),
        state: state
      }
    }).pipe(takeUntil(this.$destroy));
  }

  private _updateState(state: RskContributionState, comment?: string) {
    this.loading.push(true);
    this.apiClient.requestChunkContributionState({
      contributionId: this.contribution.id,
      body: {
        requestState: state,
        comment: comment,
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

  markRejected(comment: string) {
    this._updateState(RskContributionState.STATE_REJECTED, comment);
  }

  discard() {
    if (confirm('are you sure you want to discard saved transcript?')) {
      this.apiClient.deleteChunkContribution({ contributionId: this.contribution.id });
      this.router.navigate(['/chunk', this.chunk.id]);
    }
  }
}
