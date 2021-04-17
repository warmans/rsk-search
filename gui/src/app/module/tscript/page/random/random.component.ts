import { Component, OnInit } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { Router } from '@angular/router';
import { RskChunkStats } from '../../../../lib/api-client/models';

@Component({
  selector: 'app-random',
  templateUrl: './random.component.html',
  styleUrls: ['./random.component.scss']
})
export class RandomComponent implements OnInit {

  loading = false;

  chunkStats: RskChunkStats;

  constructor(private apiClient: SearchAPIClient, private router: Router) {
  }

  ngOnInit(): void {
    this.loading = true;
    this.apiClient.getChunkStats().subscribe((stats: RskChunkStats) => {
      this.chunkStats = stats;
      if (stats.suggestedNextChunkId) {
        this.router.navigate(['/chunk', stats.suggestedNextChunkId]);
      }
    }).add(() => {
      this.loading = false;
    });
  }

}
