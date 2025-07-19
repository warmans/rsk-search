import {Component, EventEmitter, Input} from '@angular/core';
import {
  RskAuthorContribution,
  RskAuthorContributionList,
  RskChunkedTranscriptStats,
  RskContributionState, RskShortTranscriptChange,
  RskTranscriptChange,
  RskTranscriptChangeList
} from "../../../../lib/api-client/models";
import {And, Eq, Neq} from "../../../../lib/filter-dsl/filter";
import {Bool, Str} from "../../../../lib/filter-dsl/value";
import {takeUntil} from "rxjs/operators";
import {SearchAPIClient} from "../../../../lib/api-client/services/search";

@Component({
    selector: 'app-changes',
    templateUrl: './changes.component.html',
    styleUrls: ['./changes.component.scss'],
    standalone: false
})
export class ChangesComponent {


  private unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  loading: boolean[] = [];

  pendingChanges: RskTranscriptChange[] = [];
  unapprovedPendingChanges: number = 0;

  recentContributions: RskAuthorContribution[] = [];

  activeTab: 'progress' | 'pending' | 'recent' = 'recent';

  // map of tscript_id => { 'approved' => 1, 'pending_approval' => 2 ...}
  progressMap: { [index: string]: { [index: string]: number } } = {};

  overallTotal: number = 0;
  overallComplete: number = 0;
  overallPendingApproval: number = 0;

  @Input()
  set chunks(value: RskChunkedTranscriptStats[]) {
    this.progressMap = {};

    this.overallTotal = 0;
    this.overallComplete = 0;
    this.overallPendingApproval = 0;

    value.forEach((ts: RskChunkedTranscriptStats) => {
      if (this.progressMap[ts.id] === undefined) {
        this.progressMap[ts.id] = {'total': 0, 'complete': 0, 'pending_approval': 0};
      }
      for (let chunkID in ts.chunkContributions) {
        this.progressMap[ts.id]['total']++;
        this.overallTotal++;

        ts.chunkContributions[chunkID].states.forEach((sta: RskContributionState) => {
          switch (sta) {
            case RskContributionState.STATE_APPROVED:
              this.progressMap[ts.id]['complete']++;
              this.overallComplete++;
              break;
            case RskContributionState.STATE_REQUEST_APPROVAL:
              this.progressMap[ts.id]['pending_approval']++;
              this.overallPendingApproval++;
              break;
          }
        });
      }
    });
    this._chunks = value;
    if ((this._chunks || []).length > 0) {
      this.activeTab = 'progress';
    }
  }

  get chunks(): RskChunkedTranscriptStats[] {
    return this._chunks;
  }

  private _chunks: RskChunkedTranscriptStats[] = [];

  public constructor(private apiClient: SearchAPIClient) {
  }

  public ngOnInit() {
    this.getPendingChanges();
    this.getRecentContributions();
  }

  private getRecentContributions() {
    this.loading.push(true);
    this.apiClient.listAuthorContributions({
      pageSize: 10,
      sortField: 'created_at',
      sortDirection: 'DESC'
    }).pipe(takeUntil(this.unsubscribe$)).subscribe((res: RskAuthorContributionList) => {
      this.recentContributions = res.contributions;
    }).add(() => {
      this.loading.pop();
    });
  }

  private getPendingChanges() {
    this.loading.push(true);
    this.apiClient.listTranscriptChanges({
      filter: And(
        Eq('merged', Bool(false)),
        Neq('state', Str('pending')),
        Neq('state', Str('rejected')),
      ).print()
    }).pipe(takeUntil(this.unsubscribe$)).subscribe((res: RskTranscriptChangeList) => {
      this.pendingChanges = res.changes;
      this.unapprovedPendingChanges = 0;
      (res.changes || []).forEach((ch: RskShortTranscriptChange) => {
        if (ch.state === RskContributionState.STATE_REQUEST_APPROVAL) {
          this.unapprovedPendingChanges += 1;
        }
      });
      if (this.pendingChanges.length > 0 && (this._chunks || []).length === 0) {
        this.activeTab = 'pending';
      }
    }).add(() => {
      this.loading.pop();
    });
  }

}
