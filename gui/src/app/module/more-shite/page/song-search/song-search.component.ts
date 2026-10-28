import { Component, OnDestroy, OnInit } from '@angular/core';
import { UntypedFormControl, ReactiveFormsModule } from '@angular/forms';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';
import { RskSong, RskSongList } from '../../../../lib/api-client/models';
import { ActivatedRoute, ParamMap, Router, RouterLink } from '@angular/router';
import { Like, Or } from '../../../../lib/filter-dsl/filter';
import { Str } from '../../../../lib/filter-dsl/value';
import { NgClass } from '@angular/common';

const PAGE_SIZE = 25;
const MAX_PAGINATION_LINKS = 10;
const DEFAULT_SORTING = 'title';
const SORT_OPTIONS = [
  { value: 'title', label: 'Title (A-Z)' },
  { value: 'first_occurrence', label: 'First occurrence (oldest first)' },
  { value: '-first_occurrence', label: 'First occurrence (newest first)' },
];

@Component({
  selector: 'app-song-search',
  templateUrl: './song-search.component.html',
  styleUrl: './song-search.component.scss',
  imports: [ReactiveFormsModule, RouterLink, NgClass],
})
export class SongSearchComponent implements OnInit, OnDestroy {
  searchInput: UntypedFormControl = new UntypedFormControl('');
  currentSorting: string = DEFAULT_SORTING;
  sortDropdownOpen: boolean = false;
  sortOptions = SORT_OPTIONS;
  destroy$: Subject<void> = new Subject<void>();

  public currentPage: number = 1;

  public songs: RskSong[] = [];
  public pages: number[] = [];
  public morePages: boolean = false;
  public maxPages: number = MAX_PAGINATION_LINKS;

  constructor(
    private apiClient: SearchAPIClient,
    private route: ActivatedRoute,
    private router: Router,
  ) {}

  ngOnInit(): void {
    this.route.queryParamMap.pipe(takeUntil(this.destroy$)).subscribe((route: ParamMap) => {
      this.currentPage = parseInt(route.get('page')) || 1;
      this.searchInput.setValue(route.get('q') || '', { emitEvent: false });
      this.currentSorting = (route.get('sort') || DEFAULT_SORTING).trim();
      this.updateResults();
    });
    this.searchInput.valueChanges.pipe(takeUntil(this.destroy$), debounceTime(500)).subscribe((searchTerm) => {
      this.router.navigate(['/more-shite', 'song-search'], {
        queryParams: { q: searchTerm.trim(), sort: this.currentSorting || DEFAULT_SORTING },
      });
      this.currentPage = 1;
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  get currentSortLabel(): string {
    return this.sortOptions.find((option) => option.value === this.currentSorting)?.label || SORT_OPTIONS[0].label;
  }

  setSorting(sorting: string) {
    this.sortDropdownOpen = false;
    this.router.navigate(['/more-shite', 'song-search'], {
      queryParams: { sort: sorting, page: 1 },
      queryParamsHandling: 'merge',
    });
  }

  updateResults() {
    const sorting: string = this.currentSorting || DEFAULT_SORTING;
    const sortDirection: string = sorting.startsWith('-') ? 'desc' : 'asc';
    const sortField: string = sorting.startsWith('-') ? sorting.substring(1) : sorting;
    this.apiClient
      .listSongs({
        filter:
          this.searchInput.value.trim() === ''
            ? ''
            : Or(Like('title', Str(this.searchInput.value.trim())), Like('artist', Str(this.searchInput.value.trim()))).print(),
        sortField: sortField,
        sortDirection: sortDirection,
        page: this.currentPage,
        pageSize: PAGE_SIZE,
      })
      .pipe(takeUntil(this.destroy$))
      .subscribe((res: RskSongList) => {
        this.songs = res.songs;
        let totalPages: number = Math.ceil(res.resultCount / PAGE_SIZE);
        this.pages = Array(Math.min(totalPages, MAX_PAGINATION_LINKS))
          .fill(0)
          .map((x, i) => i + 1);
        this.morePages = totalPages > MAX_PAGINATION_LINKS;
      });
  }
}
