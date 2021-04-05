import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { RsksearchAuthorRanking } from '../../../../lib/api-client/models';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-leaderboard',
  templateUrl: './leaderboard.component.html',
  styleUrls: ['./leaderboard.component.scss']
})
export class LeaderboardComponent implements OnInit, OnDestroy {

  authors: RsksearchAuthorRanking[] = [];

  showAwardHelp: boolean = false;

  private destroy$ = new EventEmitter<boolean>();

  constructor(private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {
    this.apiClient.searchServiceGetAuthorLeaderboard().pipe(takeUntil(this.destroy$)).subscribe((res) => {
      this.authors = res.authors;
    });
  }

  ngOnDestroy(): void {
    this.destroy$.emit(true);
    this.destroy$.complete();
  }

  counter(i: number) {
    return new Array(i);
  }

  nextRewardAt(i: number): string {
    return `${10 - ( i % 10)}`
  }

}
