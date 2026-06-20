import { Component, EventEmitter, OnDestroy } from '@angular/core';
import { SearchAPIClient } from 'lib/api-client/services/search';
import { takeUntil } from 'rxjs/operators';
import { RskShortTranscript, RskTranscriptList } from 'lib/api-client/models';
import { EpisodeSummaryComponent } from '../episode-summary/episode-summary.component';
import { LoadingSpinnerComponent } from '../../../shared/component/loading-spinner/loading-spinner.component';
import { NgClass } from '@angular/common';
import { Neq } from 'lib/filter-dsl/filter';
import { Null } from 'lib/filter-dsl/value';

@Component({
  selector: 'app-episode-ratings',
  templateUrl: './episode-ratings.component.html',
  styleUrls: ['./episode-ratings.component.scss'],
  imports: [LoadingSpinnerComponent, EpisodeSummaryComponent, NgClass],
})
export class EpisodeRatingsComponent implements OnDestroy {
  loading: boolean[] = [];

  transcriptList: RskShortTranscript[] = [];

  sortDirection: 'desc' | 'asc' = 'desc';

  private destroy$ = new EventEmitter<void>();

  constructor(private apiClient: SearchAPIClient) {
    this.listEpisodes();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  setSortDirection(direction: 'desc' | 'asc'): void {
    if (this.sortDirection === direction) {
      return;
    }
    this.sortDirection = direction;
    this.listEpisodes();
  }

  listEpisodes(): void {
    this.transcriptList = [];
    this.loading.push(true);
    this.apiClient
      .listTranscripts({ sortField: 'rating_score', sortDirection: this.sortDirection, pageSize: 20, filter: Neq('rating_score', Null()).print() })
      .pipe(takeUntil(this.destroy$))
      .subscribe((res: RskTranscriptList) => {
        this.transcriptList = res.episodes;
      })
      .add(() => {
        this.loading.pop();
      });
  }
}
