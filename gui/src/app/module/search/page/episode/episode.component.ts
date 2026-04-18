import { Component, EventEmitter, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Data, NavigationEnd, Router, RouterLink } from '@angular/router';
import { SearchAPIClient } from 'lib/api-client/services/search';
import { DialogType, RskArchive, RskDialog, RskTranscript, RskTranscriptChange, RskTranscriptChangeList } from 'lib/api-client/models';
import { Location, ViewportScroller, NgClass, DecimalPipe, DatePipe } from '@angular/common';
import { filter, takeUntil } from 'rxjs/operators';
import { Title } from '@angular/platform-browser';
import { SessionService } from '../../../core/service/session/session.service';
import { And, Eq, Neq } from 'lib/filter-dsl/filter';
import { Bool, Str } from 'lib/filter-dsl/value';
import { MetaService } from '../../../core/service/meta/meta.service';
import { AudioService, PlayerState, Status } from '../../../core/service/audio/audio.service';
import { Section, TranscriptComponent } from '../../../shared/component/transcript/transcript.component';
import { combineLatest } from 'rxjs';
import { parseSection } from '../../../shared/lib/fragment';
import { ClipboardService } from '../../../core/service/clipboard/clipboard.service';
import { CommunityAPIClient } from 'lib/api-client/services/community';
import { episodeIdVariations } from 'lib/util';
import { AlertService } from '../../../core/service/alert/alert.service';
import { StarRatingModule } from 'angular-star-rating';
import { SynopsesComponent } from '../../../shared/component/synopses/synopses.component';
import { MarkdownComponent } from '../../../shared/component/markdown/markdown.component';
import { AudioWaveformComponent } from '../../../shared/component/audio-waveform/audio-waveform.component';
import { LoadingOverlayComponent } from '../../../shared/component/loading-overlay/loading-overlay.component';
import { TimecodeAccuracyPipe } from '../../pipe/timecode-accuracy.pipe';

@Component({
  selector: 'app-episode',
  templateUrl: './episode.component.html',
  styleUrls: ['./episode.component.scss'],
  imports: [
    RouterLink,
    StarRatingModule,
    NgClass,
    SynopsesComponent,
    MarkdownComponent,
    TranscriptComponent,
    AudioWaveformComponent,
    LoadingOverlayComponent,
    DecimalPipe,
    DatePipe,
    TimecodeAccuracyPipe,
  ],
})
export class EpisodeComponent implements OnInit, OnDestroy {
  loading: boolean = false;

  id: string;

  shortID: string;

  scrollToID: string | null = null;

  scrollToSeconds: number | null = null;

  episode: RskTranscript;

  episodeImage: string;

  pendingChanges: RskTranscriptChange[];

  error: string;

  transcribers: string;

  quotes: RskDialog[] = [];

  songs: RskDialog[] = [];

  media: RskArchive[] = [];

  authenticated: boolean = false;

  previousEpisodeId: string;

  nextEpisodeId: string;

  audioStatus: Status;

  audioStates = PlayerState;

  selection: Section;

  unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  activeInfoPanel: 'synopsis' | 'songs' | 'quotes' | 'trivia' | 'media' = 'synopsis';

  showDownloadDialog: boolean = false;

  episodeDurationMs: number = 0;

  authorIdentifier: string = '';

  @ViewChild('transcript')
  transcriptInstance: TranscriptComponent;

  constructor(
    private router: Router,
    route: ActivatedRoute,
    private apiClient: SearchAPIClient,
    private communityApiClient: CommunityAPIClient,
    private viewportScroller: ViewportScroller,
    private titleService: Title,
    sessionService: SessionService,
    private meta: MetaService,
    private audioService: AudioService,
    private clipboard: ClipboardService,
    private alertService: AlertService,
    private location: Location,
  ) {
    route.paramMap.pipe(takeUntil(this.unsubscribe$)).subscribe((d: Data) => {
      this.loadEpisode(d.params['id']);
    });
    sessionService.onTokenChange.pipe(takeUntil(this.unsubscribe$)).subscribe((token: string): void => {
      if (token != null) {
        this.authenticated = true;
        const claims = sessionService.getClaims();
        this.authorIdentifier = `${claims.oauth_provider}:${claims.identity.name}`;
      }
    });
  }

  ngOnInit(): void {
    const fragment = this.router.url.split('#')[1] ?? null;
    this.handleFragment(fragment);

    this.router.events
      .pipe(
        filter((event) => event instanceof NavigationEnd),
        takeUntil(this.unsubscribe$),
      )
      .subscribe((event) => {
        const fragment = event.url.split('#')[1] ?? null;
        this.handleFragment(fragment);
      });

    this.audioService.status.pipe(takeUntil(this.unsubscribe$)).subscribe((sta: Status) => {
      this.audioStatus = sta;
    });
  }

  loadEpisode(id: string) {
    this.id = id;
    this.loading = true;
    this.error = undefined;

    combineLatest([
      this.apiClient.getTranscript({ epid: this.id }),
      this.meta.getMeta(),
      this.communityApiClient.listArchive({ episodeIds: episodeIdVariations(this.id) }),
    ])
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(
        ([ep, metadata, media]) => {
          this.episode = ep;
          this.shortID = ep.shortId;
          this.titleService.setTitle(ep.id);
          this.transcribers = ep.contributors.join(', ');
          this.episodeImage = ep.metadata['cover_art_url'] ? ep.metadata['cover_art_url'] : `/assets/cover/${ep.publication}-s${ep.series}-lg.jpeg`;
          this.episodeDurationMs = parseInt(ep.metadata['duration_ms']);
          this.media = media.items ?? [];

          this.quotes = [];
          this.songs = [];
          ep.transcript.forEach((r: RskDialog) => {
            if (r.notable) {
              this.quotes.push(r);
            }
            if (r.type === DialogType.SONG) {
              this.songs.push(r);
            }
          });

          const curIndex = (metadata.episodeShortIds || []).findIndex((v) => v == ep.shortId);
          if (curIndex === -1) {
            console.error(`failed to find episode in metadata ${ep.shortId}`);
          }

          this.previousEpisodeId = this.nextEpisodeId = null;
          if (curIndex > 0 && (metadata.episodeShortIds || []).length > 0) {
            this.previousEpisodeId = `ep-${metadata.episodeShortIds[curIndex - 1]}`;
          }
          if (curIndex < (metadata.episodeShortIds || []).length - 1) {
            this.nextEpisodeId = `ep-${metadata.episodeShortIds[curIndex + 1]}`;
          }

          const availableInfoPanels = [];
          if (this.episode?.synopses && this.episode?.synopses.length > 0) {
            availableInfoPanels.push('synopsis');
          }
          if (this.quotes && this.episode?.synopses.length > 0) {
            availableInfoPanels.push('quotes');
          }
          if (this.songs && this.songs.length > 0) {
            availableInfoPanels.push('songs');
          }
          if ((this.episode.trivia || []).length > 0) {
            availableInfoPanels.push('trivia');
          }
          this.activeInfoPanel = availableInfoPanels.length > 0 ? availableInfoPanels[0] : undefined;
          this.selection = parseSection(this.scrollToID, this.episode?.transcript || []);
        },
        (err) => {
          this.error = 'Failed to fetch episode';
          console.error(err);
        },
      )
      .add(() => (this.loading = false));

    this.apiClient
      .listTranscriptChanges({
        filter: And(Eq('epid', Str(this.id)), Eq('merged', Bool(false)), Neq('state', Str('pending')), Neq('state', Str('rejected'))).print(),
      })
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((ep: RskTranscriptChangeList) => {
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

  onAudioTimestamp(offsetMs: number) {
    this.audioService.setAudioSrcFromEpisodeName(this.episode.shortId, this.episode.name);
    this.audioService.seekAudio(offsetMs / 1000);
    this.audioService.playAudio();
  }

  playAudio() {
    this.audioService.setAudioSrcFromEpisodeName(this.episode.shortId, this.episode.name);
    this.audioService.playAudio();
  }

  pauseAudio() {
    this.audioService.pauseAudio();
  }

  selectSection(sel: Section) {
    this.router.navigate([], { fragment: `pos-${sel.startPos}-${sel.endPos}` });
  }

  clearSelection() {
    const includesFragment = this.location.path(true).includes('#');
    if (includesFragment) {
      this.location.go(this.location.path(false));
      this.handleFragment(null);
    }
  }

  copySelection() {
    this.clipboard.copyTextToClipboard(
      (this.episode.transcript.slice(this.selection.startPos - 1, this.selection.endPos) || [])
        .map((d: RskDialog): string => (d.actor ? `**${d.actor}:** ${d.content}` : `*${d.content}*`))
        .join('\n'),
    );
  }

  toggleDownloadDialog() {
    this.audioService.pauseAudio();
    this.showDownloadDialog = !this.showDownloadDialog;
  }

  protected readonly Math = Math;

  scrollToSelection() {
    this.transcriptInstance.scrollToAnchor();
  }

  protected readonly Array = Array;

  rateEpisode(userScore: { rating: number }) {
    if (!userScore || !this.authenticated) {
      return;
    }
    this.apiClient.setTranscriptRatingScore({ epid: this.episode.shortId, body: { score: userScore.rating } }).subscribe(() => {
      this.alertService.success('Rating submitted');
      this.episode.ratings.scores[this.authorIdentifier] = userScore.rating;
    });
  }

  private handleFragment(fragment: string | null): void {
    if (!fragment) {
      this.scrollToID = null;
      this.scrollToSeconds = null;
      return;
    }
    if (fragment.startsWith('pos-')) {
      this.scrollToID = fragment;
      this.selection = parseSection(fragment, this.episode?.transcript || []);
    }
    if (fragment.startsWith('sec-')) {
      this.scrollToSeconds = parseInt(fragment.replace('sec-', ''));
    }
  }
}
