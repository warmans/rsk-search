import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { RskAuthorRank, RskAuthorRankList } from 'src/app/lib/api-client/models';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { FormControl } from '@angular/forms';
import { Like } from 'src/app/lib/filter-dsl/filter';
import { Str } from 'src/app/lib/filter-dsl/value';

@Component({
  selector: 'app-rankings',
  templateUrl: './rankings.component.html',
  styleUrls: ['./rankings.component.scss']
})
export class RankingsComponent implements OnInit, OnDestroy {

  private destroy$ = new EventEmitter<boolean>();

  ranking: RskAuthorRank[] = [];

  searchInput: FormControl = new FormControl('');

  loading: boolean = false;
  showMoreInfo: boolean;

  constructor(private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {
    this.refreshRankings();
    this.searchInput.valueChanges.pipe(takeUntil(this.destroy$), debounceTime(100)).subscribe((val) => {
      this.refreshRankings(val);
    });
  }

  refreshRankings(username?: string) {
    this.loading = true;
    this.apiClient.listAuthorRanks(
      { filter: (username || '').trim() ? Like('author_name', Str(username.trim())).print() : '' },
    ).pipe(takeUntil(this.destroy$)).subscribe((res: RskAuthorRankList) => {
      this.ranking = res.rankings;
    }).add(() => {
      this.loading = false;
    });
  }

  ngOnDestroy(): void {
    this.destroy$.emit(true);
    this.destroy$.complete();
  }
}
