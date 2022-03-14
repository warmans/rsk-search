import { Component, EventEmitter, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Data } from '@angular/router';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import {
  RskDialog,
  RskTranscript,
  RskTranscriptChange,
  RskTranscriptChangeList
} from '../../../../lib/api-client/models';
import { ViewportScroller } from '@angular/common';
import { takeUntil } from 'rxjs/operators';
import { Title } from '@angular/platform-browser';
import { AudioPlayerComponent } from '../../../shared/component/audio-player/audio-player.component';
import { SessionService } from '../../../core/service/session/session.service';
import { And, Eq, Neq } from '../../../../lib/filter-dsl/filter';
import { Bool, Str } from '../../../../lib/filter-dsl/value';

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

  pendingChanges: RskTranscriptChange[];

  error: string;

  audioLink: string;

  transcribers: string;

  quotes: RskDialog[] = [];

  authenticated: boolean = false;

  unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  @ViewChild('audioPlayer')
  audioPlayer: AudioPlayerComponent;

  constructor(
    private route: ActivatedRoute,
    private apiClient: SearchAPIClient,
    private viewportScroller: ViewportScroller,
    private titleService: Title,
    private sessionService: SessionService,
  ) {
    route.paramMap.subscribe((d: Data) => {
      this.id = d.params['id'];
      this.shortID = this.id.replace(/ep\-/, '');
      this.audioLink = `https://storage.googleapis.com/scrimpton-raw-audio/${this.shortID}.mp3`;
    });
    route.fragment.subscribe((f) => {
      this.scrollToID = f;
    });
    sessionService.onTokenChange.pipe(takeUntil(this.unsubscribe$)).subscribe((token: string): void => {
      if (token != null) {
        this.authenticated = true;
      }
    });
  }

  ngOnInit(): void {
    this.loading = true;
    this.error = undefined;
    this.apiClient.getTranscript({ epid: this.id }).pipe(takeUntil(this.unsubscribe$)).subscribe(
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

    this.apiClient.listTranscriptChanges({
      filter: And(Eq('epid', Str(this.id)), Eq('merged', Bool(false)), Neq('state', Str('pending')), Neq('state', Str('rejected'))).print()
    }).pipe(takeUntil(this.unsubscribe$)).subscribe((ep: RskTranscriptChangeList) => {
      this.pendingChanges = ep.changes;
    });
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

  onAudioTimestamp(ts: number) {
    this.audioPlayer.seek(ts, true);
  }
}
