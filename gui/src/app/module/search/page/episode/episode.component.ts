import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Data } from '@angular/router';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { RskDialog, RskTranscript } from '../../../../lib/api-client/models';
import { ViewportScroller } from '@angular/common';
import { takeUntil } from 'rxjs/operators';
import { Title } from '@angular/platform-browser';

@Component({
  selector: 'app-episode',
  templateUrl: './episode.component.html',
  styleUrls: ['./episode.component.scss']
})
export class EpisodeComponent implements OnInit, OnDestroy {

  loading: boolean = false;

  id: string;

  shortID: string;

  scrollToID: string;

  episode: RskTranscript;

  error: string;

  audioLink: string;

  transcribers: string;

  quotes: RskDialog[] = [];

  unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(
    private route: ActivatedRoute,
    private apiClient: SearchAPIClient,
    private viewportScroller: ViewportScroller,
    private titleService: Title
  ) {
    route.paramMap.subscribe((d: Data) => {
      this.id = d.params['id'];
      this.shortID = this.id.replace(/ep\-/, '');
      this.audioLink = `https://storage.googleapis.com/scrimpton-raw-audio/${this.shortID}.mp3`;
    });
    route.fragment.subscribe((f) => {
      this.scrollToID = f;
    });
  }

  ngOnInit(): void {
    this.loading = true;
    this.error = undefined;
    this.apiClient.getTranscript({ id: this.id }).pipe(takeUntil(this.unsubscribe$)).subscribe(
      (ep: RskTranscript) => {
        this.episode = ep;
        this.titleService.setTitle(ep.id);
        this.transcribers = ep.contributors.join(', ');
        ep.transcript.forEach((r: RskDialog) => {
          if (r.notable) {
            this.quotes.push(r);
          }
        });
      },
      (err) => {
        this.error = 'Failed to fetch episode';
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
