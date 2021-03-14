import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { RskSearchResultList, RsksearchShortEpisode } from '../../../../lib/api-client/models';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss']
})
export class SearchComponent implements OnInit, OnDestroy {

  loading: boolean = false;

  episodeList: RsksearchShortEpisode[] = [];

  result: RskSearchResultList;
  pages: number[] = [];
  currentPage: number;
  morePages: boolean = false;

  private unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private apiClient: SearchAPIClient, private route: ActivatedRoute) {
    route.queryParamMap.pipe(takeUntil(this.unsubscribe$)).subscribe((params: ParamMap) => {
      this.currentPage = parseInt(params.get('page'), 10) || 0;
      if (params.get('q') === null || params.get('q').trim() == '') {
        this.result = null;
        return;
      }
      this.executeQuery(params.get('q'), this.currentPage);
    });
  }

  ngOnInit(): void {
    this.listEpisodes();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next(true);
    this.unsubscribe$.complete();
  }

  listEpisodes() {
    this.apiClient.searchServiceListEpisodes().subscribe((res) => {
      this.episodeList = res.episodes;
    });
  }

  executeQuery(value: string, page: number) {
    this.result = undefined;
    this.loading = true;
    this.apiClient.searchServiceSearch({ query: value, page: page }).subscribe((res) => {
      this.result = res;
      let totalPages = Math.ceil(res.resultCount / 15);
      this.pages = Array(Math.min(totalPages, 10)).fill(0).map((x, i) => i);
      this.morePages = totalPages > 10;
    }).add(() => {
      this.loading = false;
    });
  }
}
