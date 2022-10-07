import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Data } from '@angular/router';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { RskTranscript } from '../../../../lib/api-client/models';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { AudioService } from '../../../core/service/audio/audio.service';

@Component({
  selector: 'app-transcript-section',
  templateUrl: './transcript-section.component.html',
  styleUrls: ['./transcript-section.component.scss']
})
export class TranscriptSectionComponent implements OnInit, OnDestroy {

  unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  audioLink: string;
  epid: string;
  startLine: number;
  endLine: number;
  episode: RskTranscript;
  loading: boolean;
  error: string;

  constructor(private route: ActivatedRoute, private apiClient: SearchAPIClient, private audioService: AudioService) {

  }

  ngOnInit(): void {
    this.route.queryParamMap.pipe(takeUntil(this.unsubscribe$), debounceTime(1000)).subscribe((d: Data) => {
      this.epid = d.params['epid'];
      this.startLine = parseInt(d.params['start'], 10) || 0;
      this.endLine = parseInt(d.params['end'], 10) || this.startLine+1;
      this.loadEpisode();
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next(true);
    this.unsubscribe$.complete();
  }

  loadEpisode() {
    this.loading = true;
    this.error = undefined;
    this.apiClient.getTranscript({ epid: this.epid }).pipe(takeUntil(this.unsubscribe$)).subscribe((ep: RskTranscript) => {
        this.episode = ep;
        this.audioLink = ep.audioUri;

        // find the audio range
        let startSecond, endSecond: number;
        ep.transcript.forEach((line, k) => {
          if (k === this.startLine) {
            startSecond = parseInt(line.offsetSec);
          }
          if (k === this.endLine) {
            endSecond = parseInt(line.offsetSec);
          }
        })
        this.audioService.setAudioSrc(ep.id, ep.name, ep.audioUri, true, startSecond, endSecond);
      },
      (err) => {
        this.error = 'Failed to fetch episode';
      }).add(() => this.loading = false);
  }
}
