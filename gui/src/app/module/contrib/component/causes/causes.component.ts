import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { takeUntil } from 'rxjs/operators';
import { RskDonationStats, RskRecipientStats } from 'src/app/lib/api-client/models';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-causes',
  templateUrl: './causes.component.html',
  styleUrls: ['./causes.component.scss']
})
export class CausesComponent implements OnInit, OnDestroy {

  private destroy$ = new EventEmitter<boolean>();

  stats: RskRecipientStats[] = [];

  totalPoints: number;

  totalUSD: number;

  showMoreInfo: boolean = false;

  pointsForReward: number = environment.contributionAwardThreshold

  constructor(private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {

    this.totalPoints = 0;
    this.totalUSD = 0;

    this.apiClient.getDonationStats().pipe(takeUntil(this.destroy$)).subscribe((res: RskDonationStats) => {
      this.stats = res.stats;
      this.stats.forEach((stat) => {
        this.totalPoints += stat.pointsSpent;
        this.totalUSD += stat.donatedAmountUsd;
      })
    });
  }

  ngOnDestroy(): void {
    this.destroy$.emit(true);
    this.destroy$.complete();
  }

}
