import { Component, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { RskSearchResultList } from '../../../../lib/api-client/models';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss']
})
export class SearchComponent implements OnInit {

  result: RskSearchResultList;
  pages: number[] = [];
  currentPage: number;
  morePages: boolean = false;

  constructor(private apiClient: SearchAPIClient, private route: ActivatedRoute) {
    route.queryParamMap.subscribe((params) => {
      this.currentPage = parseInt(params.get('page'), 10) || 0;

      if (params.get('q') === null || params.get('q').trim() == '') {
        this.result = null;
        return;
      }
      this.executeQuery(params.get('q'), this.currentPage);
    });
  }

  ngOnInit(): void {

  }

  executeQuery(value: string, page: number) {
    console.log('searching...', value);
    this.apiClient.searchServiceSearch({ query: value, page: page }).subscribe((res) => {
      console.log(res);
      this.result = res;

      let totalPages = Math.ceil(res.resultCount / 15);
      this.pages = Array(Math.min(totalPages, 10)).fill(0).map((x, i)=> i);
      this.morePages = totalPages > 10;
    });
  }
}
