import { Component, OnInit } from '@angular/core';
import { RsksearchChunkStats } from '../../../../lib/api-client/models';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { Router } from '@angular/router';

@Component({
  selector: 'app-random',
  templateUrl: './random.component.html',
  styleUrls: ['./random.component.scss']
})
export class RandomComponent implements OnInit {

  loading = false;

  chunkStats: RsksearchChunkStats;

  constructor(private apiClient: SearchAPIClient, private router: Router) {
  }

  ngOnInit(): void {
    this.loading = true;
    this.apiClient.searchServiceGetTscriptChunkStats().subscribe((stats: RsksearchChunkStats) => {
      this.chunkStats = stats;
      if (stats.suggestedNextChunkId) {
        this.router.navigate(['/chunk', stats.suggestedNextChunkId]);
      }
    }).add(() => {
      this.loading = false;
    });
  }

}
