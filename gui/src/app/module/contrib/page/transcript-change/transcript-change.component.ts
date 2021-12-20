import { Component, EventEmitter, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { takeUntil } from 'rxjs/operators';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { Title } from '@angular/platform-browser';
import { SessionService } from '../../../core/service/session/session.service';
import { AlertService } from '../../../core/service/alert/alert.service';
import {
  RskContributionState,
  RskTranscript,
  RskTranscriptChange,
  RskTranscriptChangeDiff
} from '../../../../lib/api-client/models';
import { TranscriberComponent } from '../../../shared/component/transcriber/transcriber.component';
import { Observable } from 'rxjs';

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

  updatedTranscript: string;

  authenticated: boolean = false;
  userCanEdit: boolean = true;
  userIsOwner: boolean = true;
  userIsApprover: boolean = false;
  cStates = RskContributionState;
  lastUpdateTimestamp: Date;
  unifiedDiff: string;

  loading: boolean[] = [];

  $destroy: EventEmitter<boolean> = new EventEmitter<boolean>();

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

    route.paramMap.pipe(takeUntil(this.$destroy)).subscribe((d: Data) => {

      this.epID = d.params['epid'];

      this.apiClient.getTranscript({
        epid: this.epID,
        withRaw: true
      }).pipe(takeUntil(this.$destroy)).subscribe((res: RskTranscript) => {
        this.transcript = res;
        this.initialTranscript = res.rawTranscript;
      });

      if (d.params['change_id']) {
        this.apiClient.getTranscriptChange({ id: d.params['change_id'] }).pipe(takeUntil(this.$destroy)).subscribe((res: RskTranscriptChange) => {
          this.change = res;
          this.checkUserCanEdit();
          this.getDiff();

          this.initialTranscript = this.change.transcript;
          this.updatedTranscript = this.change.transcript;

          this.userIsOwner = this.sessionService.getClaims()?.author_id === res.author.id || this.sessionService.getClaims().approver;
          this.userIsApprover = this.sessionService.getClaims().approver;
        });
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
  }

  handleSave(transcript: string): void {
    this.updatedTranscript = transcript;
    if (this.change && this.userCanEdit) {
      this.update();
    }
  }

  create() {
    if (!this.change) {
      this.loading.push(true);
      this.apiClient.createTranscriptChange({
        epid: this.transcript.id,
        body: { epid: this.transcript.id, transcript: this.updatedTranscript }
      }).pipe(takeUntil(this.$destroy)).subscribe((res: RskTranscriptChange) => {
        this.transcriber.clearBackup();
        this.alertService.success('Created', 'Draft change was created. It will now be auto-saved on change.');
        this.router.navigate(['/ep', this.transcript.id, 'change', res.id]);
      }).add(() => this.loading.shift());
    }
  }

  update() {
    this._update(this.change.state).subscribe((res: RskTranscriptChange) => {
      this.change = res;
      this.checkUserCanEdit();
      this.lastUpdateTimestamp = new Date();
      this.getDiff();
    });
  }

  private _update(state: RskContributionState): Observable<RskTranscriptChange> {
    return this.apiClient.updateTranscriptChange({
      id: this.change.id,
      body: {
        id: this.change.id,
        transcript: this.updatedTranscript,
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
    const isAuthorOrApprover = this.sessionService.getClaims()?.author_id === this.change.author.id || this.sessionService.getClaims().approver;
    this.userCanEdit = this.change.state === RskContributionState.STATE_PENDING && isAuthorOrApprover;
    console.log(this.change.state === RskContributionState.STATE_PENDING, isAuthorOrApprover)
  }

  markComplete() {
    this._updateState(RskContributionState.STATE_REQUEST_APPROVAL);
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

  getDiff() {
    this.loading.push(true);
    this.apiClient.getTranscriptChangeDiff({
      id: this.change.id,
    }).pipe(takeUntil(this.$destroy)).subscribe((res: RskTranscriptChangeDiff) => {
      this.unifiedDiff = res.diff;
    }).add(() => this.loading.shift());
  }

}
