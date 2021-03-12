import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { SearchService } from '../../../module/core/service/search/search.service';

@Component({
  selector: 'app-root',
  templateUrl: './root.component.html',
  styleUrls: ['./root.component.scss']
})
export class RootComponent {
  title = 'RSK DB';

  constructor(private router: Router, private searchService: SearchService) {

  }

  executeSearch(query: string) {
    this.router.navigate(['/search'], { queryParams: { q: query } });
  }

}
