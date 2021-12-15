import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { takeUntil } from 'rxjs/operators';
import { Title } from '@angular/platform-browser';
import { RskChunkStats, RskSearchResultList } from '../../../../lib/api-client/models';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss']
})
export class SearchComponent implements OnInit, OnDestroy {

  loading: boolean[] = [];

  result: RskSearchResultList;
  pages: number[] = [];
  currentPage: number;
  morePages: boolean = false;

  chunkStats: RskChunkStats;

  private unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private apiClient: SearchAPIClient, private route: ActivatedRoute, private titleService: Title) {

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
    this.titleService.setTitle('Scrimpton Search');
    this.getChunkStats();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next(true);
    this.unsubscribe$.complete();
  }

  getChunkStats() {
    this.loading.push(true);
    this.apiClient.getChunkStats().subscribe((stats: RskChunkStats) => {
      this.chunkStats = stats;
    }).add(() => {
      this.loading.pop();
    });
  }

  executeQuery(value: string, page: number) {
    this.result = undefined;
    this.loading.push(true);
    this.apiClient.search({
      query: value,
      page: page
    }).pipe(
      takeUntil(this.unsubscribe$),
    ).subscribe((res: RskSearchResultList) => {
      this.result = res;
      let totalPages = Math.ceil(res.resultCount / 15);
      this.pages = Array(Math.min(totalPages, 10)).fill(0).map((x, i) => i);
      this.morePages = totalPages > 10;
    }).add(() => {
      this.loading.pop();
    });
  }
}
