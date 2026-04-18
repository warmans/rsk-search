import { Component, Input, inject, OnInit, DestroyRef } from '@angular/core';
import { AudioService, PlayerMode, PlayerState, Status } from '../../../core/service/audio/audio.service';
import { UntypedFormControl } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { QuotaService } from 'src/app/module/core/service/quota/quota.service';
import { RskQuotas } from 'src/app/lib/api-client/models';
import { addMonths, intervalToDuration, startOfMonth } from 'date-fns';

@Component({
  selector: 'app-audio-player-v2',
  templateUrl: './audio-player-v2.component.html',
  styleUrls: ['./audio-player-v2.component.scss'],
  standalone: false,
})
export class AudioPlayerV2Component implements OnInit {
  @Input()
  showCloseControl: boolean = true;

  audioStatus: Status;

  states = PlayerState;

  volumeControl: UntypedFormControl = new UntypedFormControl(100);

  playerProgressControl: UntypedFormControl = new UntypedFormControl(0);

  timeTillQuotaRefreshed: string = 'unknown';
  bandwidthQuotaUsedPcnt: number = 0;

  mobileMenuOpen = false;
  mobileVolumeOpen = false;

  protected readonly PlayerMode = PlayerMode;

  private audioService = inject(AudioService);
  private quotaService = inject(QuotaService);
  private destroyRef = inject(DestroyRef);

  ngOnInit(): void {
    const interval = intervalToDuration({
      start: new Date(),
      end: startOfMonth(addMonths(new Date(), 1)),
    });
    this.timeTillQuotaRefreshed = interval.days > 0 ? `${interval.days} days` : interval.hours > 0 ? `${interval.hours} hours` : `${interval.minutes} minutes`;

    this.quotaService.quotas$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((res: RskQuotas) => {
      this.bandwidthQuotaUsedPcnt = (1 - res.bandwidthRemainingMib / res.bandwidthTotalMib) * 100;
    });

    this.audioService.status.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((sta: Status) => {
      this.audioStatus = sta;
      if (!sta) {
        return;
      }
      this.playerProgressControl.setValue(sta.currentTime, { emitEvent: false });
    });

    // volume is loaded from local storage across session
    this.audioService.volumeLoaded.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((vol) => {
      this.volumeControl.setValue(vol * 100);
    });

    // persist it back if the volume is changed via our control.
    this.volumeControl.valueChanges.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((v) => {
      this.audioService.setVolume(v / 100);
    });

    this.playerProgressControl.valueChanges.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((v) => {
      this.audioService.seekAudio(v);
    });
  }

  play() {
    this.audioService.playAudio();
  }

  pause() {
    this.audioService.pauseAudio();
  }

  skipForward() {
    this.audioService.seekAudio(this.audioStatus.currentTime + 30);
  }

  closeAudio() {
    this.audioService.reset();
  }

  markAsPlayed() {
    this.audioService.markAsPlayed();
  }

  markAsUnplayed() {
    this.audioService.markAsUnplayed();
  }

  activeSleepTimer() {
    this.audioService.incrementSleepTimer(60 * 15);
  }

  deactivateSleepTimer() {
    this.audioService.clearSleepTimer();
  }
}
