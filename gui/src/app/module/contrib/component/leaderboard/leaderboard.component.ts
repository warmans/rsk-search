import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { RskAuthorLeaderboard, RskAuthorRanking } from '../../../../lib/api-client/models';
import { takeUntil } from 'rxjs/operators';
import { environment } from '../../../../../environments/environment.prod';

@Component({
  selector: 'app-leaderboard',
  templateUrl: './leaderboard.component.html',
  styleUrls: ['./leaderboard.component.scss']
})
export class LeaderboardComponent implements OnInit, OnDestroy {

  authors: RskAuthorRanking[] = [];

  showAwardHelp: boolean = false;

  awardThreshold = environment.contributionAwardThreshold;

  private destroy$ = new EventEmitter<boolean>();

  constructor(private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {
    this.apiClient.getAuthorLeaderboard().pipe(takeUntil(this.destroy$)).subscribe((res: RskAuthorLeaderboard) => {
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
    return `${environment.contributionAwardThreshold - (i % environment.contributionAwardThreshold)}`;
  }

}
