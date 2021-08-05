import { Component, EventEmitter, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { takeUntil } from 'rxjs/operators';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { Title } from '@angular/platform-browser';
import { SessionService } from '../../../core/service/session/session.service';
import { AlertService } from '../../../core/service/alert/alert.service';
import { RskContributionState, RskTranscript, RskTranscriptChange } from '../../../../lib/api-client/models';
import { TranscriberComponent } from '../../../shared/component/transcriber/transcriber.component';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-transcript-change',
  templateUrl: './transcript-change.component.html',
  styleUrls: ['./transcript-change.component.scss']
})
export class TranscriptChangeComponent implements OnInit, OnDestroy {

  epID: string;

  transcript: RskTranscript;

  change: RskTranscriptChange;

  updatedTranscript: string;

  authenticated: boolean = false;
  userCanEdit: boolean = true;
  userIsOwner: boolean = true;
  userIsApprover: boolean = false;

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
      });

      if (d.params['change_id']) {
        this.apiClient.getTranscriptChange({ id: d.params['change_id'] }).pipe(takeUntil(this.$destroy)).subscribe((res: RskTranscriptChange) => {
          this.change = res;
          this.updatedTranscript = this.change.transcript;

          this.userCanEdit = res.state === RskContributionState.STATE_PENDING;
          if (!this.sessionService.getClaims().approver) {
            this.userCanEdit = this.sessionService.getClaims()?.author_id === res.author.id;
          }
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
      //this.lastUpdateTimestamp = new Date();
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

}
