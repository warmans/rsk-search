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

  private destroy$ = new EventEmitter<boolean>();

  constructor(private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {
    this.apiClient.searchServiceGetAuthorLeaderboard().pipe(takeUntil(this.destroy$)).subscribe((res) => {
      this.authors = res.authors;
      console.log(this.authors);
    });
  }

  ngOnDestroy(): void {
    this.destroy$.emit(true);
    this.destroy$.complete();
  }

  counter(i: number) {
    return new Array(i);
  }

}
