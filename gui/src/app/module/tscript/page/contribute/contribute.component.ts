import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import {
  RsksearchContributionState,
  RsksearchTscriptList,
  RsksearchTscriptStats
} from '../../../../lib/api-client/models';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-contribute',
  templateUrl: './contribute.component.html',
  styleUrls: ['./contribute.component.scss']
})
export class ContributeComponent implements OnInit, OnDestroy {

  loading: boolean[] = [];

  tscipts: RsksearchTscriptStats[] = [];

  // map of tscript_id => { 'approved' => 1, 'pending_approval' => 2 ...}
  progressMap: { [index: string]: { [index: string]: number } } = {};

  private unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {
    this.loading.push(true);
    this.apiClient.searchServiceListTscripts().pipe(takeUntil(this.unsubscribe$)).subscribe((res: RsksearchTscriptList) => {
      this.tscipts = res.tscripts;
      this.progressMap = {};
      this.tscipts.forEach((ts) => {
        if (this.progressMap[ts.id] === undefined) {
          this.progressMap[ts.id] = { 'total': 0, 'complete': 0, 'pending_approval': 0 };
        }
        for (let chunkID in ts.chunkContributions) {
          this.progressMap[ts.id]['total']++;
          ts.chunkContributions[chunkID].states.forEach((sta) => {
            switch (sta) {
              case RsksearchContributionState.STATE_APPROVED:
                this.progressMap[ts.id]['complete']++;
                break;
              case RsksearchContributionState.STATE_REQUEST_APPROVAL:
                this.progressMap[ts.id]['pending_approval']++;
                break;
            }
          });
        }
        console.log(this.progressMap);
      });
    }).add(() => {
      this.loading.pop();
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.emit(true);
    this.unsubscribe$.complete();
  }

}
