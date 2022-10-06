import { Component, EventEmitter, Input, OnDestroy, OnInit } from '@angular/core';
import { RskShortTranscript } from '../../../../lib/api-client/models';
import { AudioService } from '../../../core/service/audio/audio.service';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-episode-summary',
  templateUrl: './episode-summary.component.html',
  styleUrls: ['./episode-summary.component.scss']
})
export class EpisodeSummaryComponent implements OnInit, OnDestroy {

  @Input()
  set episode(value: RskShortTranscript) {
    this._episode = value;
    this.episodeImage = value?.metadata['cover_art_url'] ? value?.metadata['cover_art_url'] : `/assets/cover/${value.publication}-s${value.series}.jpg`;
  }

  get episode(): RskShortTranscript {
    return this._episode;
  }

  private _episode: RskShortTranscript;

  episodeImage: string;

  played: boolean = false;

  private destroy$ = new EventEmitter<void>();

  constructor(private audioService: AudioService) {
  }

  ngOnInit(): void {
    this.audioService.audioHistoryLog.pipe(takeUntil(this.destroy$)).subscribe((played) => {
      if (!this.episode) {
        return;
      }
      console.log(played, this.episode.shortId);
      this.played = played.indexOf(this.episode.shortId) > -1
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  playEpisode(episode: RskShortTranscript) {
    this.audioService.setAudioSrc(episode.shortId, episode.audioUri);
    this.audioService.playAudio();
  }
}
