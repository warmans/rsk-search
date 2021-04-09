import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { takeUntil } from 'rxjs/operators';
import { RsksearchReward } from '../../../../lib/api-client/models';
import { Router, RoutesRecognized } from '@angular/router';

@Component({
  selector: 'app-pending-rewards',
  templateUrl: './pending-rewards.component.html',
  styleUrls: ['./pending-rewards.component.scss']
})
export class PendingRewardsComponent implements OnInit, OnDestroy {

  displayOnPage: boolean = true;

  rewards: RsksearchReward[];

  private destroy$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private apiClient: SearchAPIClient, private router: Router) {
  }

  ngOnInit(): void {
    this.router.events.subscribe((data) => {
      if (data instanceof RoutesRecognized) {
        this.displayOnPage = !data?.state?.root?.firstChild?.data?.disableRewardPopup;
        if (this.displayOnPage) {
          this.apiClient.searchServiceListPendingRewards().pipe(takeUntil(this.destroy$)).subscribe((res) => {
            this.rewards = res.rewards;
          });
        }
      }
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next(true);
    this.destroy$.complete();
  }

}
