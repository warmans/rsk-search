import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { AudioService, PlayerState, Status } from '../../../core/service/audio/audio.service';
import { FormControl } from '@angular/forms';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-audio-player-v2',
  templateUrl: './audio-player-v2.component.html',
  styleUrls: ['./audio-player-v2.component.scss']
})
export class AudioPlayerV2Component implements OnInit, OnDestroy {

  audioStatus: Status;

  states = PlayerState;

  volumeControl: FormControl = new FormControl(100);

  playerProgressControl: FormControl = new FormControl(0);

  private unsubscribe$: EventEmitter<void> = new EventEmitter<void>();

  constructor(private audioService: AudioService) {
  }

  ngOnInit(): void {

    this.audioService.status.pipe(takeUntil(this.unsubscribe$)).subscribe((sta: Status) => {
      this.audioStatus = sta;
      if (!sta){
        return;
      }
      this.playerProgressControl.setValue(sta.currentTime, {emitEvent: false});
    });

    this.volumeControl.valueChanges.pipe(takeUntil(this.unsubscribe$)).subscribe((v) => {
      this.audioService.setVolume(v/100);
    });

    this.playerProgressControl.valueChanges.pipe(takeUntil(this.unsubscribe$)).subscribe((v) => {
      this.audioService.seekAudio(v);
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
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


}
