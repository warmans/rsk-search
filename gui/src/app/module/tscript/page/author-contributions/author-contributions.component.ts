import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { Claims, SessionService } from '../../../core/service/session/session.service';
import { takeUntil } from 'rxjs/operators';
import {
  RsksearchChunkContribution,
  RsksearchChunkContributionList,
  RsksearchContributionState
} from '../../../../lib/api-client/models';
import { Router } from '@angular/router';

@Component({
  selector: 'app-author-contributions',
  templateUrl: './author-contributions.component.html',
  styleUrls: ['./author-contributions.component.scss']
})
export class AuthorContributionsComponent implements OnInit, OnDestroy {

  claims: Claims;

  contributions: RsksearchChunkContribution[];

  loading: boolean = false;

  states = RsksearchContributionState;

  private destroy$: EventEmitter<any> = new EventEmitter<any>();

  constructor(private apiClient: SearchAPIClient, private session: SessionService, private router: Router) {
  }

  ngOnInit(): void {
    this.session.onTokenChange.pipe(takeUntil(this.destroy$)).subscribe((v) => {
      if (v) {
        this.claims = this.session.getClaims();
        this.loadContributions();
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
      this.apiClient.searchServiceDiscardDraftContribution({
        chunkId: chunkId,
        contributionId: contributionId
      }).pipe(takeUntil(this.destroy$)).subscribe(() => {
        this.loadContributions();
      });
    }
  }

  loadContributions() {
    this.loading = true;
    this.apiClient.searchServiceListAuthorContributions({
      authorId: this.session.getClaims().author_id,
      page: 0
    }).pipe(takeUntil(this.destroy$)).subscribe((list: RsksearchChunkContributionList) => {
      this.contributions = list.contributions;
    }).add(() => this.loading = false);
  }

  logout() {
    this.session.destroySession();
    this.router.navigate(['/search']);
  }
}
