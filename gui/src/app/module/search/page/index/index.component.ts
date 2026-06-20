import { Component, EventEmitter, HostListener, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from 'lib/api-client/services/search';
import { ActivatedRoute, ParamMap, Router, RouterLink } from '@angular/router';
import { takeUntil } from 'rxjs/operators';
import { Title } from '@angular/platform-browser';
import {
  RskChangelog,
  RskChunkedTranscriptList,
  RskChunkedTranscriptStats,
  RskDialog,
  RskSearchResult,
  RskSearchResultList,
  RskShortTranscript,
} from 'lib/api-client/models';
import { AudioService } from '../../../core/service/audio/audio.service';
import { ClipboardService } from 'module/core/service/clipboard/clipboard.service';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { SearchStatsComponent } from '../../component/search-stats/search-stats.component';
import { TranscriptCopyDialogComponent } from '../../../shared/component/transcript-copy-dialog/transcript-copy-dialog.component';
import { TranscriptComponent } from '../../../shared/component/transcript/transcript.component';
import { NgClass } from '@angular/common';
import { MarkdownComponent } from '../../../shared/component/markdown/markdown.component';
import { EpisodeListComponent } from '../../component/episode-list/episode-list.component';
import { EpisodeRatingsComponent } from '../../component/episode-ratings/episode-ratings.component';
import { LoadingOverlayComponent } from '../../../shared/component/loading-overlay/loading-overlay.component';
import { MatchedRowPosPipe } from '../../pipe/match-row-pos.pipe';

@Component({
  selector: 'app-search-index',
  templateUrl: './index.component.html',
  styleUrls: ['./index.component.scss'],
  imports: [
    SearchStatsComponent,
    ReactiveFormsModule,
    RouterLink,
    TranscriptCopyDialogComponent,
    TranscriptComponent,
    NgClass,
    MarkdownComponent,
    EpisodeListComponent,
    EpisodeRatingsComponent,
    LoadingOverlayComponent,
    MatchedRowPosPipe,
  ],
})
export class IndexComponent implements OnInit, OnDestroy {
  loading: boolean[] = [];

  query: string;
  result: RskSearchResultList;
  results: RskSearchResult[] = [];
  currentPage: number = 0;
  currentSorting = new FormControl<string>('_score');
  morePages: boolean = false;
  private loadingNextPage: boolean = false;
  latestChangelog: RskChangelog;
  contributionsNeeded: number;
  banner: { image: string; url: string };

  private unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  activeInfoPanel: 'contribute' | 'changelog' | 'roadmap' = 'contribute';
  roadmapMarkdown: string = 'Loading...';

  activeSection: 'episodes' | 'ratings' = 'episodes';

  sectionDropdownOpen: boolean = false;

  constructor(
    private apiClient: SearchAPIClient,
    private route: ActivatedRoute,
    private titleService: Title,
    private audioService: AudioService,
    private clipboardService: ClipboardService,
    private router: Router,
  ) {
    //this.banner = {image: 'partridge-banner.png', url: 'https://discord.gg/nKABACyy6d'},
    this.banner =
      new Date().getMonth() >= 10
        ? { image: 'pilk-christmas-banner.png', url: 'https://woodymakesgames.itch.io/averypilkingtonchristmas' }
        : { image: 'partridge-banner.png', url: 'https://discord.gg/nKABACyy6d' };

    this.currentSorting.valueChanges.pipe(takeUntil(this.unsubscribe$)).subscribe((val) => {
      if (!val) {
        return;
      }
      this.router.navigate([], {
        queryParams: {
          sort: val,
        },
        queryParamsHandling: 'merge',
        skipLocationChange: false,
      });
    });

    route.queryParamMap.pipe(takeUntil(this.unsubscribe$)).subscribe((params: ParamMap) => {
      let sorting = (params.get('sort') || '').trim();
      if (sorting) {
        this.currentSorting.setValue(sorting);
      }
      this.activeSection = params.get('section') === 'ratings' ? 'ratings' : 'episodes';
      this.query = (params.get('q') || '').trim();
      if (this.query === '') {
        this.result = null;
        this.results = [];
        return;
      }
      this.executeQuery(this.query, 0, this.currentSorting.getRawValue() ?? '_score');
    });
  }

  ngOnInit(): void {
    this.titleService.setTitle('Scrimpton Search');

    this.apiClient
      .listChangelogs({ pageSize: 1 })
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((res) => {
        this.latestChangelog = (res.changelogs || []).pop();
      });

    this.apiClient
      .listChunkedTranscripts()
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((res: RskChunkedTranscriptList) => {
        this.contributionsNeeded = 0;
        (res.chunked || []).forEach((v: RskChunkedTranscriptStats) => {
          this.contributionsNeeded += v.numChunks - (v.numApprovedContributions || 0);
        });
      });

    this.apiClient
      .getRoadmap()
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((res) => {
        this.roadmapMarkdown = res.markdown;
      });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next(true);
    this.unsubscribe$.complete();
  }

  executeQuery(value: string, page: number, sort: string, append: boolean = false) {
    if (!append) {
      this.result = undefined;
      this.results = [];
    }
    this.loadingNextPage = true;
    this.loading.push(true);
    this.apiClient
      .search({
        query: value,
        page: page,
        sort: sort,
      })
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe((res: RskSearchResultList) => {
        this.result = res;
        this.results = append ? this.results.concat(res.results || []) : res.results || [];
        this.currentPage = page;
        let totalPages = Math.ceil(res.resultCount / 15);
        this.morePages = page + 1 < totalPages;
      })
      .add(() => {
        this.loading.pop();
        this.loadingNextPage = false;
      });
  }

  loadNextPage() {
    if (this.loadingNextPage || !this.morePages || !this.query) {
      return;
    }
    this.executeQuery(this.query, this.currentPage + 1, this.currentSorting.getRawValue() ?? '_score', true);
  }

  @HostListener('window:scroll')
  onWindowScroll() {
    if (this.query === '' || !this.morePages || this.loadingNextPage) {
      return;
    }
    const threshold = 600;
    const scrollPosition = window.innerHeight + window.scrollY;
    const documentHeight = document.documentElement.scrollHeight;
    if (scrollPosition >= documentHeight - threshold) {
      this.loadNextPage();
    }
  }

  onAudioTimestamp(ep: RskShortTranscript, tsMs: number) {
    this.audioService.setAudioSrcFromEpisodeName(ep.shortId, ep.name);
    this.audioService.seekAudio(tsMs / 1000);
    this.audioService.playAudio();
  }

  copyLineToClipboard(line: RskDialog) {
    this.clipboardService.copyTextToClipboard(line.content);
  }
}
