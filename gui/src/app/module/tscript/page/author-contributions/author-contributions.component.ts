import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { Claims, SessionService } from '../../../core/service/session/session.service';
import { takeUntil } from 'rxjs/operators';
import { Router } from '@angular/router';
import { Eq } from '../../../../lib/filter-dsl/filter';
import { Str } from '../../../../lib/filter-dsl/value';
import {
  RskClaimedReward,
  RskContribution,
  RskContributionList,
  RskContributionState
} from '../../../../lib/api-client/models';
import { Title } from '@angular/platform-browser';

@Component({
  selector: 'app-author-contributions',
  templateUrl: './author-contributions.component.html',
  styleUrls: ['./author-contributions.component.scss']
})
export class AuthorContributionsComponent implements OnInit, OnDestroy {

  claims: Claims;

  contributions: RskContribution[];

  rewards: RskClaimedReward[];

  loading: boolean[] = [];

  states = RskContributionState;

  private destroy$: EventEmitter<any> = new EventEmitter<any>();

  constructor(private apiClient: SearchAPIClient, private session: SessionService, private router: Router, private titleService: Title) {
    titleService.setTitle('Author Contributions');
  }

  ngOnInit(): void {
    this.session.onTokenChange.pipe(takeUntil(this.destroy$)).subscribe((v) => {
      if (v) {
        this.claims = this.session.getClaims();
        this.loadContributions();
        this.loadClaimedRewards();
      } else {
        this.claims = undefined;
      }
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next(null);
    this.destroy$.complete();
  }

  discardDraft(chunkId: string, contributionId: string): void {
    if (confirm('Really discard draft?')) {
      this.apiClient.deleteContribution({
        contributionId: contributionId
      }).pipe(takeUntil(this.destroy$)).subscribe(() => {
        this.loadContributions();
      });
    }
  }

  loadContributions() {
    this.loading.push(true);
    this.apiClient.listContributions({
      filter: Eq(`author_id`, Str(this.session.getClaims().author_id)).print(),
      sortField: `created_at`,
      sortDirection: 'desc',
    }).pipe(takeUntil(this.destroy$)).subscribe((list: RskContributionList) => {
      this.contributions = list.contributions;
    }).add(() => this.loading.pop());
  }

  loadClaimedRewards() {
    this.loading.push(true);
    this.apiClient.listClaimedRewards({}).pipe(takeUntil(this.destroy$)).subscribe((list) => {
      this.rewards = list.rewards;
    }).add(() => this.loading.pop());
  }

  logout() {
    this.session.destroySession();
    this.router.navigate(['/search']);
  }
}
