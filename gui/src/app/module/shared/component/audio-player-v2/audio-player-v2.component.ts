import { Component, OnInit } from '@angular/core';
import { AudioService, PlayerState, Status } from '../../../core/service/audio/audio.service';

@Component({
  selector: 'app-audio-player-v2',
  templateUrl: './audio-player-v2.component.html',
  styleUrls: ['./audio-player-v2.component.scss']
})
export class AudioPlayerV2Component implements OnInit {

  audioStatus: Status;

  states = PlayerState;

  constructor(private audioService: AudioService) { }

  ngOnInit(): void {

    this.audioService.status.subscribe((sta: Status) => {
      console.log(sta);
      this.audioStatus = sta;
    })

    this.audioService.setAudioSrc("xfm-S1E01", 'https://storage.googleapis.com/scrimpton-raw-audio/xfm-S1E01.mp3');
  }

  play() {

    this.audioService.playAudio();
  }

  pause() {
    this.audioService.pauseAudio();
  }

  skipForward() {
    this.audioService.seekAudio(this.audioStatus.currentTime + 30)
  }
}
