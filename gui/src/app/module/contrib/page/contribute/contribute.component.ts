import {Component, EventEmitter, OnDestroy, OnInit} from '@angular/core';
import {SearchAPIClient} from 'src/app/lib/api-client/services/search';
import {
  RskChunkedTranscriptList,
  RskChunkedTranscriptStats,
  RskContributionState,
} from 'src/app/lib/api-client/models';
import {takeUntil} from 'rxjs/operators';
import {Title} from '@angular/platform-browser';

@Component({
  selector: 'app-contribute',
  templateUrl: './contribute.component.html',
  styleUrls: ['./contribute.component.scss']
})
export class ContributeComponent implements OnInit, OnDestroy {

  loading: boolean[] = [];

  chunkedTranscripts: RskChunkedTranscriptStats[] = [];

  activeContributionsPanel: 'authors' | 'outgoing_donations' | 'incoming_donations' = 'authors';

  private unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  overallAwaitingContributions: number = 0;

  constructor(private apiClient: SearchAPIClient, private titleService: Title) {
    titleService.setTitle('Contribute');
  }

  ngOnInit(): void {
    this.loading.push(true);
    this.apiClient.listChunkedTranscripts().pipe(takeUntil(this.unsubscribe$)).subscribe((res: RskChunkedTranscriptList) => {
      this.chunkedTranscripts = res.chunked;

      let overallTotal: number = 0;
      let overallComplete: number = 0;
      let overallPendingApproval: number = 0;
      this.chunkedTranscripts.forEach((ts: RskChunkedTranscriptStats) => {
        for (let chunkID in ts.chunkContributions) {
          overallTotal++;
          ts.chunkContributions[chunkID].states.forEach((sta: RskContributionState) => {
            switch (sta) {
              case RskContributionState.STATE_APPROVED:
                overallComplete++;
                break;
              case RskContributionState.STATE_REQUEST_APPROVAL:
                overallPendingApproval++;
                break;
            }
          });
        }
      });
      this.overallAwaitingContributions = overallTotal - (overallComplete + overallPendingApproval);

    }).add(() => {
      this.loading.pop();
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.emit(true);
    this.unsubscribe$.complete();
  }

}
