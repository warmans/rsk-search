import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { Claims, SessionService } from '../../../core/service/session/session.service';
import { takeUntil } from 'rxjs/operators';
import { Router } from '@angular/router';
import { Eq } from 'src/app/lib/filter-dsl/filter';
import { Str } from 'src/app/lib/filter-dsl/value';
import {
  RskChunkContribution,
  RskChunkContributionList,
  RskClaimedReward,
  RskContributionState,
  RskShortTranscriptChange,
  RskTranscriptChangeList
} from 'src/app/lib/api-client/models';
import { Title } from '@angular/platform-browser';

@Component({
  selector: 'app-author-profile',
  templateUrl: './author-profile.component.html',
  styleUrls: ['./author-profile.component.scss']
})
export class AuthorProfile implements OnInit, OnDestroy {

  claims: Claims;

  contributions: RskChunkContribution[];

  changes: RskShortTranscriptChange[];

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
        this.loadChanges();
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
      this.apiClient.deleteChunkContribution({
        contributionId: contributionId
      }).pipe(takeUntil(this.destroy$)).subscribe(() => {
        this.loadContributions();
      });
    }
  }

  loadContributions() {
    this.loading.push(true);
    this.apiClient.listChunkContributions({
      filter: Eq(`author_id`, Str(this.session.getClaims().author_id)).print(),
      sortField: `created_at`,
      sortDirection: 'desc',
    }).pipe(takeUntil(this.destroy$)).subscribe((list: RskChunkContributionList) => {
      this.contributions = list.contributions;
    }).add(() => this.loading.pop());
  }

  loadChanges() {
    this.loading.push(true);
    this.apiClient.listTranscriptChanges({
      filter: Eq(`author_id`, Str(this.session.getClaims().author_id)).print(),
      sortField: `created_at`,
      sortDirection: 'desc',
    }).pipe(takeUntil(this.destroy$)).subscribe((list: RskTranscriptChangeList) => {
      this.changes = list.changes;
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
