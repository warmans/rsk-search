import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { takeUntil } from 'rxjs/operators';
import { RskAuthorRank, RskAuthorRankList } from '../../../../lib/api-client/models';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';

@Component({
  selector: 'app-rankings',
  templateUrl: './rankings.component.html',
  styleUrls: ['./rankings.component.scss']
})
export class RankingsComponent implements OnInit, OnDestroy {

  private destroy$ = new EventEmitter<boolean>();

  ranking: RskAuthorRank[] = [];

  constructor(private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {
    this.apiClient.listAuthorRanks({}).pipe(takeUntil(this.destroy$)).subscribe((res: RskAuthorRankList) => {
      this.ranking = res.rankings;
    });
  }

  ngOnDestroy(): void {
    this.destroy$.emit(true);
    this.destroy$.complete();
  }
}
