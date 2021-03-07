import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Data } from '@angular/router';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { RsksearchEpisode } from '../../../../lib/api-client/models';
import { ViewportScroller } from '@angular/common';

@Component({
  selector: 'app-episode',
  templateUrl: './episode.component.html',
  styleUrls: ['./episode.component.scss']
})
export class EpisodeComponent implements OnInit {

  id: string;

  scrollToID: string;

  episode: RsksearchEpisode;

  constructor(
    private route: ActivatedRoute,
    private apiClient: SearchAPIClient,
    private viewportScroller: ViewportScroller,
  ) {
    route.paramMap.subscribe((d: Data) => {
      this.id = d.params['id'];
    });
    route.fragment.subscribe((f) => {
      this.scrollToID = f;
    });
  }

  ngOnInit(): void {
    this.apiClient.searchServiceGetEpisode({ id: this.id }).subscribe((ep: RsksearchEpisode) => {
      this.episode = ep;
    });
  }

  query(field: string, value: string): string {
    return `${field} = "${value}"`;
  }

  scrollToTop() {
    this.viewportScroller.scrollToPosition([0, 0]);
  }
}
