import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { takeUntil } from 'rxjs/operators';
import { RskReward } from 'src/app/lib/api-client/models';
import { Router, RoutesRecognized } from '@angular/router';
import { SessionService } from '../../../core/service/session/session.service';

@Component({
  selector: 'app-pending-rewards',
  templateUrl: './pending-rewards.component.html',
  styleUrls: ['./pending-rewards.component.scss']
})
export class PendingRewardsComponent implements OnInit, OnDestroy {

  displayOnPage: boolean = true;

  rewards: RskReward[];
  rewardIcons: string[];

  prizeIcons: string[] = [
    '/assets/prizes/ladder49-80px.png',
    '/assets/prizes/children-of-the-corn-80px.png',
    '/assets/prizes/executive-decision-80px.png',
    '/assets/prizes/stigmata-80px.png',
    '/assets/prizes/scotland-rocks-80px.png',
    '/assets/prizes/best-air-guitar-80px.png',
  ];

  private destroy$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private apiClient: SearchAPIClient, private router: Router, private session: SessionService) {
  }

  ngOnInit(): void {
    this.router.events.subscribe((data) => {
      if (data instanceof RoutesRecognized) {
        this.displayOnPage = !data?.state?.root?.firstChild?.data?.disableRewardPopup;
        if (this.displayOnPage && this.session.getToken()) {
          this.apiClient.listPendingRewards().pipe(takeUntil(this.destroy$)).subscribe((res) => {
            this.rewards = res.rewards;
            this.rewardIcons = [];
            this.rewards.forEach(() => {
              this.rewardIcons.push(this.randomPrize());
            })
          });
        }
      }
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next(true);
    this.destroy$.complete();
  }

  randomPrize(): string {
    return this.prizeIcons[this.randomInt(0, this.prizeIcons.length)];
  }

  randomInt(min: number, max: number): number {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min) + min);
  }
}
