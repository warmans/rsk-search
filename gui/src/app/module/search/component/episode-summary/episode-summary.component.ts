import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  Input,
  OnDestroy,
  OnInit
} from '@angular/core';
import {RskPublicationType, RskShortTranscript} from '../../../../lib/api-client/models';
import {AudioService, Status} from '../../../core/service/audio/audio.service';
import {takeUntil} from 'rxjs/operators';
import {SessionService} from 'src/app/module/core/service/session/session.service';

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
    this.episodeIdIsMeaningful = value.publicationType === RskPublicationType.PUBLICATION_TYPE_RADIO ||
      value.publicationType === RskPublicationType.PUBLICATION_TYPE_PODCAST ||
      value.publicationType === RskPublicationType.PUBLICATION_TYPE_TV
  }

  get episode(): RskShortTranscript {
    return this._episode;
  }

  loggedIn: boolean = false;

  private _episode: RskShortTranscript;

  episodeImage: string;

  episodeIdIsMeaningful: boolean = true;

  played: boolean = false;

  playing: boolean = false;

  private destroy$ = new EventEmitter<void>();

  constructor(private audioService: AudioService, private cdr: ChangeDetectorRef, private session: SessionService,) {
    session.onTokenChange.pipe(takeUntil(this.destroy$)).subscribe((token) => {
      if (token) {
        this.loggedIn = !!this.session.getClaims();
      } else {
        this.loggedIn = false;
      }
    });
  }

  ngOnInit(): void {
    this.audioService.status.pipe(takeUntil(this.destroy$)).subscribe((status: Status) => {
      const playing = (status.audioID === this.episode.shortId);
      if (playing !== this.playing) {
        this.playing = playing;
        this.cdr.detectChanges();
      }
    });

    this.audioService.audioHistoryLog.pipe(takeUntil(this.destroy$)).subscribe((played: string[]) => {
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
    if (!episode.media.audio || episode.media.audioRestricted) {
      return;
    }
    this.audioService.setAudioSrc(episode.shortId, episode.name);
    this.audioService.toggleAudio();
  }
}
