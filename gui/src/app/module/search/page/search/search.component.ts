import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { takeUntil } from 'rxjs/operators';
import { Title } from '@angular/platform-browser';
import { RskChangelog, RskChunkedTranscriptList, RskChunkedTranscriptStats, RskSearchResultList, RskShortTranscript } from 'src/app/lib/api-client/models';
import { AudioService } from '../../../core/service/audio/audio.service';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss']
})
export class SearchComponent implements OnInit, OnDestroy {

  loading: boolean[] = [];

  query: string;
  result: RskSearchResultList;
  pages: number[] = [];
  currentPage: number;
  morePages: boolean = false;
  latestChangelog: RskChangelog;
  contributionsNeeded: number;

  private unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  activeInfoPanel: 'contribute' | 'changelog' = 'contribute';

  constructor(
    private apiClient: SearchAPIClient,
    private route: ActivatedRoute,
    private titleService: Title,
    private audioService: AudioService) {

    route.queryParamMap.pipe(takeUntil(this.unsubscribe$)).subscribe((params: ParamMap) => {
      this.currentPage = parseInt(params.get('page'), 10) || 0;
      this.query = (params.get('q') || '').trim();
      if (this.query === '') {
        this.result = null;
        return;
      }
      this.executeQuery(this.query, this.currentPage);
    });
  }

  ngOnInit(): void {
    this.titleService.setTitle('Scrimpton Search');

    this.apiClient.listChangelogs({ pageSize: 1 }).pipe(takeUntil(this.unsubscribe$)).subscribe((res) => {
      this.latestChangelog = (res.changelogs || []).pop();
    });

    this.apiClient.listChunkedTranscripts().pipe(takeUntil(this.unsubscribe$)).subscribe((res: RskChunkedTranscriptList) => {
      this.contributionsNeeded = 0;
      (res.chunked || []).forEach((v: RskChunkedTranscriptStats) => {
        this.contributionsNeeded += v.numChunks - ((v.numApprovedContributions || 0) + (v.numPendingContributions || 0));
      });
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next(true);
    this.unsubscribe$.complete();
  }

  executeQuery(value: string, page: number) {
    this.result = undefined;
    this.loading.push(true);
    this.apiClient.search({
      query: value,
      page: page
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

  onAudioTimestamp(ep: RskShortTranscript, ts: number) {
    this.audioService.setAudioSrc(ep.shortId, ep.name, ep.audioUri);
    this.audioService.seekAudio(ts);
    this.audioService.playAudio();
  }
}
