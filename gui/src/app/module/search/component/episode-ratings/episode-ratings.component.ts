import { Component, EventEmitter, HostListener, OnDestroy } from '@angular/core';
import { SearchAPIClient } from 'lib/api-client/services/search';
import { takeUntil } from 'rxjs/operators';
import { RskShortTranscript, RskTranscriptList } from 'lib/api-client/models';
import { EpisodeSummaryComponent } from '../episode-summary/episode-summary.component';
import { LoadingSpinnerComponent } from '../../../shared/component/loading-spinner/loading-spinner.component';
import { NgClass, DecimalPipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { Neq } from 'lib/filter-dsl/filter';
import { Null } from 'lib/filter-dsl/value';

interface RatingsGridColumn {
  key: string;
  publication: string;
  series: number;
  label: string;
  firstReleaseTime: number;
  groupStart: boolean;
  groupEnd: boolean;
}

interface RatingsGridGroup {
  publication: string;
  colspan: number;
  firstReleaseTime: number;
}

interface RatingsGrid {
  groups: RatingsGridGroup[];
  columns: RatingsGridColumn[];
  episodes: number[];
  cells: { [columnKey: string]: { [episode: number]: RskShortTranscript } };
}

@Component({
  selector: 'app-episode-ratings',
  templateUrl: './episode-ratings.component.html',
  styleUrls: ['./episode-ratings.component.scss'],
  imports: [LoadingSpinnerComponent, EpisodeSummaryComponent, NgClass, RouterLink, DecimalPipe],
})
export class EpisodeRatingsComponent implements OnDestroy {
  loading: boolean[] = [];

  viewMode: 'list' | 'grid' = 'grid';

  transcriptList: RskShortTranscript[] = [];

  sortDirection: 'desc' | 'asc' = 'desc';

  grid: RatingsGrid = { groups: [], columns: [], episodes: [], cells: {} };

  private gridLoaded: boolean = false;

  private readonly listPageSize: number = 20;
  private listPage: number = 0;
  morePages: boolean = false;
  private loadingNextPage: boolean = false;

  private destroy$ = new EventEmitter<void>();

  constructor(private apiClient: SearchAPIClient) {
    this.loadGrid();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  setViewMode(mode: 'list' | 'grid'): void {
    if (this.viewMode === mode) {
      return;
    }
    this.viewMode = mode;
    if (mode === 'grid' && !this.gridLoaded) {
      this.loadGrid();
    } else if (mode === 'list') {
      this.listEpisodes();
    }
  }

  setSortDirection(direction: 'desc' | 'asc'): void {
    if (this.viewMode === 'list' && this.sortDirection === direction) {
      return;
    }
    this.viewMode = 'list';
    this.sortDirection = direction;
    this.listEpisodes();
  }

  listEpisodes(append: boolean = false): void {
    if (!append) {
      this.transcriptList = [];
      this.listPage = 0;
      this.morePages = false;
    }
    this.loadingNextPage = true;
    this.loading.push(true);
    this.apiClient
      .listTranscripts({
        sortField: 'rating_score',
        sortDirection: this.sortDirection,
        page: this.listPage,
        pageSize: this.listPageSize,
        filter: Neq('rating_score', Null()).print(),
      })
      .pipe(takeUntil(this.destroy$))
      .subscribe((res: RskTranscriptList) => {
        const episodes = res.episodes || [];
        this.transcriptList = append ? this.transcriptList.concat(episodes) : episodes;
        this.morePages = episodes.length >= this.listPageSize;
      })
      .add(() => {
        this.loading.pop();
        this.loadingNextPage = false;
      });
  }

  loadNextPage(): void {
    if (this.loadingNextPage || !this.morePages || this.viewMode !== 'list') {
      return;
    }
    this.listPage += 1;
    this.listEpisodes(true);
  }

  @HostListener('window:scroll')
  onWindowScroll(): void {
    if (this.viewMode !== 'list' || !this.morePages || this.loadingNextPage) {
      return;
    }
    const threshold = 600;
    const scrollPosition = window.innerHeight + window.scrollY;
    const documentHeight = document.documentElement.scrollHeight;
    if (scrollPosition >= documentHeight - threshold) {
      this.loadNextPage();
    }
  }

  loadGrid(): void {
    this.loading.push(true);
    this.apiClient
      .listTranscripts({ sortField: 'rating_score', sortDirection: 'desc', filter: Neq('rating_score', Null()).print() })
      .pipe(takeUntil(this.destroy$))
      .subscribe((res: RskTranscriptList) => {
        this.grid = this.buildGrid(res.episodes || []);
        this.gridLoaded = true;
      })
      .add(() => {
        this.loading.pop();
      });
  }

  ratingColor(rating: number): string {
    if (rating === undefined || rating === null) {
      return 'transparent';
    }
    // map rating (0-5) to a hue from red (0) to green (120)
    const clamped = Math.max(0, Math.min(5, rating));
    const hue = (clamped / 5) * 120;
    return `hsl(${hue}, 65%, 45%)`;
  }

  private buildGrid(episodes: RskShortTranscript[]): RatingsGrid {
    const columnsMap: { [key: string]: RatingsGridColumn } = {};
    const cells: { [columnKey: string]: { [episode: number]: RskShortTranscript } } = {};
    const episodeNumbers = new Set<number>();

    for (const ep of episodes) {
      if (ep.ratingScore === undefined || ep.ratingScore === null || ep.episode === undefined || ep.episode === null) {
        continue;
      }
      const publication = ep.publication || 'unknown';
      const series = ep.series || 0;
      const key = `${publication}|${series}`;
      const releaseTime = ep.releaseDate ? new Date(ep.releaseDate).getTime() : Number.MAX_SAFE_INTEGER;
      if (!columnsMap[key]) {
        columnsMap[key] = {
          key: key,
          publication: publication,
          series: series,
          label: `S${series}`,
          firstReleaseTime: releaseTime,
          groupStart: false,
          groupEnd: false,
        };
      } else if (!Number.isNaN(releaseTime) && releaseTime < columnsMap[key].firstReleaseTime) {
        columnsMap[key].firstReleaseTime = releaseTime;
      }
      if (!cells[key]) {
        cells[key] = {};
      }
      cells[key][ep.episode] = ep;
      episodeNumbers.add(ep.episode);
    }

    const columns = Object.values(columnsMap).sort((a, b) => {
      if (a.firstReleaseTime !== b.firstReleaseTime) {
        return a.firstReleaseTime - b.firstReleaseTime;
      }
      if (a.publication !== b.publication) {
        return a.publication.localeCompare(b.publication);
      }
      return a.series - b.series;
    });

    const groups: RatingsGridGroup[] = [];
    for (const col of columns) {
      const last = groups[groups.length - 1];
      if (last && last.publication === col.publication) {
        last.colspan += 1;
        if (col.firstReleaseTime < last.firstReleaseTime) {
          last.firstReleaseTime = col.firstReleaseTime;
        }
      } else {
        groups.push({ publication: col.publication, colspan: 1, firstReleaseTime: col.firstReleaseTime });
      }
    }

    for (let i = 0; i < columns.length; i++) {
      columns[i].groupStart = i === 0 || columns[i - 1].publication !== columns[i].publication;
      columns[i].groupEnd = i === columns.length - 1 || columns[i + 1].publication !== columns[i].publication;
    }

    const episodesSorted = Array.from(episodeNumbers).sort((a, b) => a - b);

    return { groups: groups, columns: columns, episodes: episodesSorted, cells: cells };
  }
}
