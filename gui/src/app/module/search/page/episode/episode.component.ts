import { Component, EventEmitter, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { RskDialog, RskMetadata, RskTranscript, RskTranscriptChange, RskTranscriptChangeList } from '../../../../lib/api-client/models';
import { ViewportScroller } from '@angular/common';
import { takeUntil } from 'rxjs/operators';
import { Title } from '@angular/platform-browser';
import { AudioPlayerComponent } from '../../../shared/component/audio-player/audio-player.component';
import { SessionService } from '../../../core/service/session/session.service';
import { And, Eq, Neq } from '../../../../lib/filter-dsl/filter';
import { Bool, Str } from '../../../../lib/filter-dsl/value';
import { MetaService } from '../../../core/service/meta/meta.service';
import { AudioService } from '../../../core/service/audio/audio.service';

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

  previousEpisodeId: string;
  nextEpisodeId: string;

  unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  @ViewChild('audioPlayer')
  audioPlayer: AudioPlayerComponent;

  constructor(
    private route: ActivatedRoute,
    private apiClient: SearchAPIClient,
    private viewportScroller: ViewportScroller,
    private titleService: Title,
    private sessionService: SessionService,
    private meta: MetaService,
    private audioService: AudioService,
  ) {
    route.paramMap.subscribe((d: Data) => {
      this.loadEpisode(d.params['id']);
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
  }

  loadEpisode(id: string) {

    this.id = id;
    this.shortID = id.replace(/ep\-/, '');
    this.audioLink = `https://storage.googleapis.com/scrimpton-raw-audio/${this.shortID}.mp3`;

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
        this.meta.getMeta().pipe(takeUntil(this.unsubscribe$)).subscribe((res: RskMetadata) => {
          const curIndex = (res.episodeShortIDs || []).findIndex((v) => v == ep.shortId);
          if (curIndex === -1) {
            console.error(`failed to find episode in metadata ${ep.shortId}`);
          }

          this.previousEpisodeId = this.nextEpisodeId = null;
          if (curIndex > 0 && (res.episodeShortIDs || []).length > 0) {
            this.previousEpisodeId = `ep-${res.episodeShortIDs[curIndex - 1]}`;
          }
          if (curIndex < ((res.episodeShortIDs || []).length - 1)) {
            this.nextEpisodeId = `ep-${res.episodeShortIDs[curIndex + 1]}`;
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

  playAudio() {
    this.audioService.setAudioSrc(this.episode.shortId, this.episode.audioUri);
    this.audioService.playAudio();
  }
}
