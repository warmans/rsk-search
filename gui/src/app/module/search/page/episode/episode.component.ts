import {Component, EventEmitter, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute, Data, Router} from '@angular/router';
import {SearchAPIClient} from 'src/app/lib/api-client/services/search';
import {
  DialogType,
  RskDialog,
  RskTranscript,
  RskTranscriptChange,
  RskTranscriptChangeList
} from 'src/app/lib/api-client/models';
import {ViewportScroller} from '@angular/common';
import {takeUntil} from 'rxjs/operators';
import {Title} from '@angular/platform-browser';
import {SessionService} from '../../../core/service/session/session.service';
import {And, Eq, Neq} from 'src/app/lib/filter-dsl/filter';
import {Bool, Str} from 'src/app/lib/filter-dsl/value';
import {MetaService} from '../../../core/service/meta/meta.service';
import {AudioService, PlayerState, Status} from '../../../core/service/audio/audio.service';
import {Section} from '../../../shared/component/transcript/transcript.component';
import {combineLatest} from 'rxjs';
import {parseSection} from "../../../shared/lib/fragment";
import {ClipboardService} from "../../../core/service/clipboard/clipboard.service";

@Component({
  selector: 'app-episode',
  templateUrl: './episode.component.html',
  styleUrls: ['./episode.component.scss']
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

  audioLink: string;

  transcribers: string;

  quotes: RskDialog[] = [];

  songs: RskDialog[] = [];

  authenticated: boolean = false;

  previousEpisodeId: string;

  nextEpisodeId: string;

  audioStatus: Status;

  audioStates = PlayerState;

  selection: Section;

  unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  activeInfoPanel: 'synopsis' | 'songs' | 'quotes' = 'synopsis';

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private apiClient: SearchAPIClient,
    private viewportScroller: ViewportScroller,
    private titleService: Title,
    private sessionService: SessionService,
    private meta: MetaService,
    private audioService: AudioService,
    private clipboard: ClipboardService,
  ) {
    route.paramMap.pipe(takeUntil(this.unsubscribe$)).subscribe((d: Data) => {
      this.loadEpisode(d.params['id']);
    });
    route.fragment.pipe(takeUntil(this.unsubscribe$)).subscribe((f) => {
      if (!f) {
        this.scrollToID = undefined;
        this.scrollToSeconds = undefined;
        return;
      }
      if (f.startsWith('pos-')) {
        this.scrollToID = f;
        this.selection = parseSection(f);
      }
      if (f.startsWith('sec-')) {
        this.scrollToSeconds = parseInt(f.replace('sec-', ''));
      }
    });
    sessionService.onTokenChange.pipe(takeUntil(this.unsubscribe$)).subscribe((token: string): void => {
      if (token != null) {
        this.authenticated = true;
      }
    });
  }

  ngOnInit(): void {
    this.audioService.status.pipe(takeUntil(this.unsubscribe$)).subscribe((sta: Status) => {
      this.audioStatus = sta;
    });
  }

  loadEpisode(id: string) {

    this.id = id;
    this.loading = true;
    this.error = undefined;

    combineLatest([
      this.apiClient.getTranscript({epid: this.id}),
      this.meta.getMeta(),
    ]).pipe(takeUntil(this.unsubscribe$)).subscribe(
      ([ep, metadata]) => {

        this.episode = ep;
        this.shortID = ep.shortId;
        this.titleService.setTitle(ep.id);
        this.transcribers = ep.contributors.join(', ');
        this.episodeImage = ep.metadata['cover_art_url'] ? ep.metadata['cover_art_url'] : `/assets/cover/${ep.publication}-s${ep.series}-lg.jpeg`;
        this.audioLink = ep.audioUri;

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
        if (curIndex < ((metadata.episodeShortIds || []).length - 1)) {
          this.nextEpisodeId = `ep-${metadata.episodeShortIds[curIndex + 1]}`;
        }

        const availableInfoPanels = [];
        if (this.episode?.synopses && this.episode?.synopses.length > 0) {
          availableInfoPanels.push("synopsis");
        }
        if (this.quotes && this.episode?.synopses.length > 0) {
          availableInfoPanels.push("quotes");
        }
        if (this.songs && this.songs.length > 0) {
          availableInfoPanels.push("songs");
        }
        this.activeInfoPanel = (availableInfoPanels.length > 0) ? availableInfoPanels[0] : undefined;
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

  onAudioTimestamp(offsetMs: number) {
    this.audioService.setAudioSrc(this.episode.shortId, this.episode.name, this.episode.audioUri);
    this.audioService.seekAudio(offsetMs/1000);
    this.audioService.playAudio();
  }

  playAudio() {
    this.audioService.setAudioSrc(this.episode.shortId, this.episode.name, this.episode.audioUri);
    this.audioService.playAudio();
  }

  pauseAudio() {
    this.audioService.pauseAudio();
  }

  selectSection(sel: Section) {
    this.router.navigate([], {fragment: `pos-${sel.startPos}-${sel.endPos}`});
  }

  clearSelection() {
    this.router.navigate([], {});
  }

  copySelection() {
    this.clipboard.copyTextToClipboard(
      (this.episode.transcript.slice(this.selection.startPos - 1, this.selection.endPos) || [])
        .map((d: RskDialog): string => d.actor ? `**${d.actor}:** ${d.content}` : `*${d.content}*`)
        .join("\n")
    );
  }
}
