import { Component, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { RskSearchResultList } from '../../../../lib/api-client/models';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss']
})
export class SearchComponent implements OnInit {

  result: RskSearchResultList;

  constructor(private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {
  }

  executeQuery(value: string) {
    console.log("searching...", value);
    this.apiClient.searchServiceSearch({ query: value }).subscribe((res) => {
      console.log(res);
      this.result = res;
    });
  }
}
