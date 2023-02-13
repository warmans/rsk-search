import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { RskChunkedTranscriptList, RskChunkedTranscriptStats, RskContributionState, RskTranscriptChange, } from 'src/app/lib/api-client/models';
import { takeUntil } from 'rxjs/operators';
import { Title } from '@angular/platform-browser';
import { And, Eq, Neq } from 'src/app/lib/filter-dsl/filter';
import { Bool, Str } from 'src/app/lib/filter-dsl/value';

@Component({
  selector: 'app-contribute',
  templateUrl: './contribute.component.html',
  styleUrls: ['./contribute.component.scss']
})
export class ContributeComponent implements OnInit, OnDestroy {

  loading: boolean[] = [];

  chunkedTranscripts: RskChunkedTranscriptStats[] = [];
  transcriptChanges: RskTranscriptChange[] = [];

  // map of tscript_id => { 'approved' => 1, 'pending_approval' => 2 ...}
  progressMap: { [index: string]: { [index: string]: number } } = {};

  overallTotal: number = 0;
  overallComplete: number = 0;
  overallPendingApproval: number = 0;
  overallAwaitingContributions: number = 0;

  activeContributionsPanel: 'authors' | 'outgoing_donations' | 'incoming_donations' = 'authors';

  private unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private apiClient: SearchAPIClient, private titleService: Title) {
    titleService.setTitle('Contribute');
  }

  ngOnInit(): void {
    this.loading.push(true);
    this.apiClient.listChunkedTranscripts().pipe(takeUntil(this.unsubscribe$)).subscribe((res: RskChunkedTranscriptList) => {
      this.chunkedTranscripts = res.chunked;

      this.progressMap = {};

      this.overallTotal = 0;
      this.overallComplete = 0;
      this.overallPendingApproval = 0;

      this.chunkedTranscripts.forEach((ts) => {
        if (this.progressMap[ts.id] === undefined) {
          this.progressMap[ts.id] = { 'total': 0, 'complete': 0, 'pending_approval': 0 };
        }
        for (let chunkID in ts.chunkContributions) {
          this.progressMap[ts.id]['total']++;
          this.overallTotal++;

          ts.chunkContributions[chunkID].states.forEach((sta) => {
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

      this.overallAwaitingContributions = this.overallTotal - (this.overallComplete + this.overallPendingApproval);

    }).add(() => {
      this.loading.pop();
    });

    this.loading.push(true);
    this.apiClient.listTranscriptChanges({
      filter: And(
        Eq('merged', Bool(false)),
        Neq('state', Str('pending')),
        Neq('state', Str('rejected')),
      ).print()
    }).pipe(takeUntil(this.unsubscribe$)).subscribe((res) => {
      this.transcriptChanges = res.changes;
    }).add(() => {
      this.loading.pop();
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.emit(true);
    this.unsubscribe$.complete();
  }

}
