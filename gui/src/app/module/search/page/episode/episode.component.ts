import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Data } from '@angular/router';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { RsksearchEpisode } from '../../../../lib/api-client/models';
import { ViewportScroller } from '@angular/common';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-episode',
  templateUrl: './episode.component.html',
  styleUrls: ['./episode.component.scss']
})
export class EpisodeComponent implements OnInit, OnDestroy {

  loading: boolean = false;

  id: string;

  scrollToID: string;

  episode: RsksearchEpisode;

  error: string;

  unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

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
    this.loading = true;
    this.error = undefined;
    this.apiClient.searchServiceGetEpisode({ id: this.id }).pipe(takeUntil(this.unsubscribe$)).subscribe(
      (ep: RsksearchEpisode) => {
        this.episode = ep;
      },
      (err) => {
        this.error = "Failed to fetch episode";
      }).add(() => this.loading = false);
  }

  query(field: string, value: string): string {
    return `${field} = "${value}"`;
  }

  scrollToTop() {
    this.viewportScroller.scrollToPosition([0, 0]);
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next(true);
    this.unsubscribe$.complete();
  }
}
