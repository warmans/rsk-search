import {Component, EventEmitter, OnDestroy, OnInit} from '@angular/core';
import {SearchAPIClient} from 'src/app/lib/api-client/services/search';
import {ActivatedRoute, ParamMap, Router} from '@angular/router';
import {takeUntil} from 'rxjs/operators';
import {Title} from '@angular/platform-browser';
import {
    RskChangelog,
    RskChunkedTranscriptList,
    RskChunkedTranscriptStats,
    RskDialog,
    RskSearchResultList,
    RskShortTranscript
} from 'src/app/lib/api-client/models';
import {AudioService} from '../../../core/service/audio/audio.service';
import {ClipboardService} from 'src/app/module/core/service/clipboard/clipboard.service';
import {FormControl} from "@angular/forms";

@Component({
    selector: 'app-search',
    templateUrl: './search.component.html',
    styleUrls: ['./search.component.scss'],
    standalone: false
})
export class SearchComponent implements OnInit, OnDestroy {

  loading: boolean[] = [];

  query: string;
  result: RskSearchResultList;
  pages: number[] = [];
  currentPage: number;
  currentSorting = new FormControl<string>("_score");
  morePages: boolean = false;
  latestChangelog: RskChangelog;
  contributionsNeeded: number;
  banner: {image: string, url : string};

  private unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  activeInfoPanel: 'contribute' | 'changelog' | 'roadmap'  = 'contribute';
  roadmapMarkdown: string = 'Loading...';

  constructor(
    private apiClient: SearchAPIClient,
    private route: ActivatedRoute,
    private titleService: Title,
    private audioService: AudioService,
    private clipboardService: ClipboardService,
    private router: Router) {

    this.banner = {image: 'partridge-banner.png', url: 'https://discord.gg/nKABACyy6d'};

    this.currentSorting.valueChanges.pipe(takeUntil(this.unsubscribe$)).subscribe((val) => {
      if (!val) {
        return;
      }
      this.router.navigate([], {
        queryParams: {
          sort: val
        },
        queryParamsHandling: 'merge',
        skipLocationChange: false,
      });
    });

    route.queryParamMap.pipe(takeUntil(this.unsubscribe$)).subscribe((params: ParamMap) => {
      this.currentPage = parseInt(params.get('page'), 10) || 0;
      let sorting = (params.get('sort') || '').trim();
      if (sorting) {
        this.currentSorting.setValue(sorting);
      }
      this.query = (params.get('q') || '').trim();
      if (this.query === '') {
        this.result = null;
        return;
      }
      this.executeQuery(
        this.query,
        this.currentPage,
        this.currentSorting.getRawValue() ?? '_score',
      );
    });
  }

  ngOnInit(): void {
    this.titleService.setTitle('Scrimpton Search');

    this.apiClient.listChangelogs({pageSize: 1}).pipe(takeUntil(this.unsubscribe$)).subscribe((res) => {
      this.latestChangelog = (res.changelogs || []).pop();
    });

    this.apiClient.listChunkedTranscripts().pipe(takeUntil(this.unsubscribe$)).subscribe((res: RskChunkedTranscriptList) => {
      this.contributionsNeeded = 0;
      (res.chunked || []).forEach((v: RskChunkedTranscriptStats) => {
        this.contributionsNeeded += v.numChunks - (v.numApprovedContributions || 0);
      });
    });

    this.apiClient.getRoadmap().pipe(takeUntil(this.unsubscribe$)).subscribe((res) => {
      this.roadmapMarkdown = res.markdown;
    })
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next(true);
    this.unsubscribe$.complete();
  }

  executeQuery(value: string, page: number, sort: string) {
    this.result = undefined;
    this.loading.push(true);
    this.apiClient.search({
      query: value,
      page: page,
      sort: sort,
    }).pipe(
      takeUntil(this.unsubscribe$),
    ).subscribe((res: RskSearchResultList) => {
      this.result = res;
      let totalPages = Math.ceil(res.resultCount / 15);
      this.pages = Array(Math.min(totalPages, 10)).fill(0).map((x, i) => i);
      this.morePages = totalPages > 10;
    }).add(() => {
      this.loading.pop();
    });
  }

  onAudioTimestamp(ep: RskShortTranscript, tsMs: number) {
    this.audioService.setAudioSrcFromEpisodeName(ep.shortId, ep.name);
    this.audioService.seekAudio(tsMs/1000);
    this.audioService.playAudio();
  }

  copyLineToClipboard(line: RskDialog) {
    this.clipboardService.copyTextToClipboard(line.content);
  }
}
