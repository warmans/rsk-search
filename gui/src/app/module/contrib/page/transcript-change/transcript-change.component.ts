import { Component, EventEmitter, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { takeUntil } from 'rxjs/operators';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { Title } from '@angular/platform-browser';
import { SessionService } from '../../../core/service/session/session.service';
import { AlertService } from '../../../core/service/alert/alert.service';
import { RskContributionState, RskTranscript, RskTranscriptChange, RskTranscriptChangeDiff } from '../../../../lib/api-client/models';
import { TranscriberComponent } from '../../../shared/component/transcriber/transcriber.component';
import { Observable } from 'rxjs';
import { FormControl } from '@angular/forms';

const DISMISS_HELP_KEY: string = 'contribute.change.help.hide';

@Component({
  selector: 'app-transcript-change',
  templateUrl: './transcript-change.component.html',
  styleUrls: ['./transcript-change.component.scss']
})
export class TranscriptChangeComponent implements OnInit, OnDestroy {

  epID: string;

  initialTranscript: string;

  transcript: RskTranscript;

  change: RskTranscriptChange;

  approvalPoints: FormControl = new FormControl(0.2);

  versionMismatchError = false;
  readOnly: boolean = true;
  authenticated: boolean = false;
  userCanEdit: boolean = true;
  userIsOwner: boolean = true;
  userIsApprover: boolean = false;
  cStates = RskContributionState;
  lastUpdateTimestamp: Date;
  unifiedDiff: string;
  instructionsHidden: boolean = localStorage.getItem(DISMISS_HELP_KEY) === 'true';

  loading: boolean[] = [];

  $destroy: EventEmitter<void> = new EventEmitter<void>();

  @ViewChild('transcriber')
  transcriber: TranscriberComponent;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private apiClient: SearchAPIClient,
    private titleService: Title,
    private sessionService: SessionService,
    private alertService: AlertService,
  ) {
    titleService.setTitle('Contribute');

    // don't bother prompting for login etc. if the intent is just to read the change.
    this.readOnly = route.snapshot.queryParamMap.get('readonly') === '1';

    route.paramMap.pipe(takeUntil(this.$destroy)).subscribe((d: Data) => {

      this.epID = d.params['epid'];

      this.loading.push(true);
      this.apiClient.getTranscript({
        epid: this.epID,
        withRaw: true
      }).pipe(takeUntil(this.$destroy)).subscribe((res: RskTranscript) => {
        this.transcript = res;
        if (!d.params['change_id']) {
          this.initialTranscript = res.rawTranscript;
        } else {
          this.apiClient.getTranscriptChange({ id: d.params['change_id'] }).pipe(takeUntil(this.$destroy)).subscribe((res: RskTranscriptChange) => {
            this.change = res;
            this.checkUserCanEdit();

            this.initialTranscript = this.change.transcript;

            this.versionMismatchError = (this.change?.transcriptVersion !== this.transcript?.version);
            this.userIsOwner = this.sessionService.getClaims()?.author_id === res.author.id || this.sessionService.getClaims()?.approver;
            this.userIsApprover = this.sessionService.getClaims()?.approver;
          });
        }
      }).add(() => this.loading.shift());
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
  }

  handleSave(transcript: string): void {
    // there is some kind of race condition that means sometimes it attempts to save an empty transcript.
    // quick-fix: just ignore updates with a falsy transcript.
    if (this.change && this.userCanEdit && transcript) {
      this.update(() => {
      });
    }
  }

  create() {
    if (!this.change) {
      this.loading.push(true);
      this.apiClient.createTranscriptChange({
        epid: this.transcript.id,
        body: { epid: this.transcript.id, transcript: this.transcriber.getContentSnapshot(), transcriptVersion: this.transcript?.version || 'NONE' }
      }).pipe(takeUntil(this.$destroy)).subscribe((res: RskTranscriptChange) => {
        this.initialTranscript = res.transcript;
        this.transcriber.clearBackup();
        this.alertService.success('Created', 'Draft change was created. It will now be auto-saved on change.');
        this.router.navigate(['/ep', this.transcript.id, 'change', res.id]);
      }).add(() => this.loading.shift());
    }
  }

  discardChange() {
    if (confirm('Really discard change?')) {
      this.loading.push(true);
      this.apiClient.deleteTranscriptChange({ id: this.change.id }).pipe(takeUntil(this.$destroy)).subscribe((res) => {
        this.router.navigate(['/ep', this.transcript.id, 'change']);
      }).add(() => this.loading.shift());
    }
  }

  update(after: () => void) {
    this._update(this.change.state).pipe(takeUntil(this.$destroy)).subscribe((res: RskTranscriptChange) => {
      this.change = res;
      this.checkUserCanEdit();
      this.lastUpdateTimestamp = new Date();
      after();
    });
  }

  private _update(state: RskContributionState): Observable<RskTranscriptChange> {
    return this.apiClient.updateTranscriptChange({
      id: this.change.id,
      body: {
        id: this.change.id,
        transcript: this.transcriber.getContentSnapshot(),
        state: state
      }
    }).pipe(takeUntil(this.$destroy));
  }

  private _updateState(state: RskContributionState) {
    this.loading.push(true);
    this.apiClient.requestTranscriptChangeState({
      id: this.change.id,
      body: {
        id: this.change.id,
        state: state,
        pointsOnApprove: this.approvalPoints.value,
      }
    }).pipe(takeUntil(this.$destroy)).subscribe((res: RskTranscriptChange) => {
      this.change.state = state;
      this.checkUserCanEdit();

      switch (this.change.state) {
        case RskContributionState.STATE_PENDING:
          this.alertService.success('Retracted', 'Change is now back in the pending state. It will not be reviewed until is is re-submitted.');
          return;
        case RskContributionState.STATE_APPROVED:
          this.alertService.success('Approved', 'Change was approved.');
          return;
        case RskContributionState.STATE_REQUEST_APPROVAL:
          this.alertService.success('Submitted', 'Change is now awaiting manual approval by an approver. This usually takes around 24 hours.');
          return;
        case RskContributionState.STATE_REJECTED:
          this.alertService.success('Rejected', 'Change was rejected.');
          return;
      }
    }).add(() => this.loading.shift());
  }

  checkUserCanEdit() {
    if (this.readOnly) {
      // don't even check if they can edit if the intent is to read only.
      this.userCanEdit = false;
    } else {
      const isAuthorOrApprover = this.sessionService.getClaims()?.author_id === this.change.author.id || this.sessionService.getClaims().approver;
      this.userCanEdit = this.change.state === RskContributionState.STATE_PENDING && isAuthorOrApprover;
    }
  }

  markComplete() {
    this._updateState(RskContributionState.STATE_REQUEST_APPROVAL);
  }

  markIncomplete() {
    this._updateState(RskContributionState.STATE_PENDING);
    // remove readonly param since it's now been moved into a writable state
    this.router.navigate([], {
      queryParams: {
        'readonly': null,
      },
      queryParamsHandling: 'merge'
    })
  }

  markApproved() {
    this._updateState(RskContributionState.STATE_APPROVED);
  }

  markRejected() {
    this._updateState(RskContributionState.STATE_REJECTED);
  }

  getDiff() {
    this.loading.push(true);
    this.apiClient.getTranscriptChangeDiff({
      id: this.change.id,
    }).pipe(takeUntil(this.$destroy)).subscribe((res: RskTranscriptChangeDiff) => {
      this.unifiedDiff = res.diff;
    }).add(() => this.loading.shift());
  }

  checkReloadDiff(v: string) {
    if (v === 'diff') {
      if (this.change.state === RskContributionState.STATE_PENDING) {
        // always update before loading the diff, to ensure it is accurate.
        this.update(() => {
          this.getDiff();
        });
      } else {
        this.getDiff();
      }
    }
  }

  hideInstructions() {
    this.instructionsHidden = true;
    localStorage.setItem(DISMISS_HELP_KEY, "true")
  }

  undoHideInstructions() {
    this.instructionsHidden = false;
    localStorage.removeItem(DISMISS_HELP_KEY)
  }
}
