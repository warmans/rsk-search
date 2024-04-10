import {Component, OnDestroy, OnInit} from '@angular/core';
import {UntypedFormControl} from "@angular/forms";
import {SearchAPIClient} from "../../../../lib/api-client/services/search";
import {debounceTime, takeUntil} from "rxjs/operators";
import {Subject} from "rxjs";
import {RskSong, RskSongList} from "../../../../lib/api-client/models";
import {ActivatedRoute, ParamMap, Router} from "@angular/router";
import {Like, Or} from "../../../../lib/filter-dsl/filter";
import {Str} from "../../../../lib/filter-dsl/value";

const PAGE_SIZE = 25;
const MAX_PAGINATION_LINKS = 10;

@Component({
  selector: 'app-song-search',
  standalone: false,
  templateUrl: './song-search.component.html',
  styleUrl: './song-search.component.scss'
})
export class SongSearchComponent implements OnInit, OnDestroy {

  searchInput: UntypedFormControl = new UntypedFormControl('');
  destroy$: Subject<void> = new Subject<void>();

  public currentPage: number = 1;

  public songs: RskSong[] = [];
  public pages: number[] = [];
  public morePages: boolean = false;
  public maxPages: number = MAX_PAGINATION_LINKS;

  constructor(private apiClient: SearchAPIClient, private route: ActivatedRoute, private router: Router) {
  }

  ngOnInit(): void {
    this.route.queryParamMap.pipe(takeUntil(this.destroy$)).subscribe((route: ParamMap) => {
      this.currentPage = parseInt(route.get("page")) || 1;
      this.searchInput.setValue(route.get("q") || "", {emitEvent: false});
      this.updateResults();
    });
    this.searchInput.valueChanges.pipe(takeUntil(this.destroy$), debounceTime(500)).subscribe((searchTerm) => {
      this.router.navigate(['/more-shite', "song-search"], {queryParams: {q: searchTerm.trim()}});
      this.currentPage = 1;
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  updateResults() {
    this.apiClient.listSongs({
      filter: this.searchInput.value.trim() === "" ? "" : Or(Like("title", Str(this.searchInput.value.trim())), Like("artist", Str(this.searchInput.value.trim()))).print(),
      page: this.currentPage,
      pageSize: PAGE_SIZE,
    }).pipe(takeUntil(this.destroy$)).subscribe((res: RskSongList) => {
      this.songs = res.songs;
      let totalPages: number = Math.ceil(res.resultCount / PAGE_SIZE);
      this.pages = Array(Math.min(totalPages, MAX_PAGINATION_LINKS)).fill(0).map((x, i) => i + 1);
      this.morePages = totalPages > MAX_PAGINATION_LINKS;
    })
  }
}
