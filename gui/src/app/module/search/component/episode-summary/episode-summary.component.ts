import { ChangeDetectionStrategy, ChangeDetectorRef, Component, EventEmitter, Input, OnDestroy, OnInit } from '@angular/core';
import { RskShortTranscript } from '../../../../lib/api-client/models';
import { AudioService, Status } from '../../../core/service/audio/audio.service';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-episode-summary',
  templateUrl: './episode-summary.component.html',
  styleUrls: ['./episode-summary.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
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

  playing: boolean = false;

  private destroy$ = new EventEmitter<void>();

  constructor(private audioService: AudioService, private cdr: ChangeDetectorRef) {
  }

  ngOnInit(): void {
    this.audioService.status.pipe(takeUntil(this.destroy$)).subscribe((status: Status) => {
      const playing = (status.audioID === this.episode.shortId);
      if (playing !== this.playing) {
        this.playing = playing;
        this.cdr.detectChanges();
      }
    });

    this.audioService.audioHistoryLog.pipe(takeUntil(this.destroy$)).subscribe((played) => {
      if (!this.episode) {
        return;
      }
      this.played = played.indexOf(this.episode.shortId) > -1;
      this.cdr.detectChanges();
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  toggleEpisode(episode: RskShortTranscript) {
    this.audioService.setAudioSrc(episode.shortId, episode.name, episode.audioUri);
    this.audioService.toggleAudio();
  }
}
