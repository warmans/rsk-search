import { AfterViewInit, Component, ElementRef, Input, ViewChild } from '@angular/core';

@Component({
  selector: 'app-audio-player',
  templateUrl: './audio-player.component.html',
  styleUrls: ['./audio-player.component.scss']
})
export class AudioPlayerComponent implements AfterViewInit {


  @Input()
  public src: string;

  @Input()
  public autoplay: boolean = false;

  @Input()
  public showStateLabel: boolean = false;

  @Input()
  public volume: number = 1.0; /* 1.0 is loudest */

  @Input()
  set playbackRate(value: number) {
    this._playbackRate = value;
    this.updatePlayer();
  }
  get playbackRate(): number {
    return this._playbackRate;
  }
  private _playbackRate: number = 1.0;

  @ViewChild('audioElement', { static: false })
  public audioPlayerEl: ElementRef;

  private audio: HTMLMediaElement;

  public constructor() {
  }

  public ngAfterViewInit() {
    if (!this.audioPlayerEl?.nativeElement) {
      return;
    }
    this.audio = this.audioPlayerEl.nativeElement;
    this.updatePlayer();
  }

  public updatePlayer() {
    if (this.audio) {
      this.audio.volume = this.volume;
      this.audio.autoplay = this.autoplay;
      this.audio.defaultPlaybackRate = this._playbackRate;
      this.audio.playbackRate = this._playbackRate;
    }
  }

  public pause(withOffset?: number): void {
    if (this.audio) {
      if (withOffset !== undefined) {
        this.audio.currentTime = this.audio.currentTime + withOffset > 0 ? this.audio.currentTime + withOffset : 0;
      }
      this.audio.pause();
    }
  }

  public play(withOffset?: number): void {
    console.log(this.audio, this.audio.readyState);
    if (this.audio && this.audio.readyState >= 1) {
        if (withOffset !== undefined) {
            this.audio.currentTime = this.audio.currentTime + withOffset > 0 ? this.audio.currentTime + withOffset : 0;
        }
        this.audio.play();
    }
  }

  public toggle(withOffset?: number) {
    if (this.playing()) {
      this.pause();
      return
    }
    this.play(withOffset);
  }

  public seek(time: number, andPlay?: boolean) {
    if (this.audio) {
      this.audio.currentTime = time;
    }
    if (andPlay) {
      this.play();
    }
  }

  public playing(): boolean {
    return !!(this.audio.currentTime > 0 && !this.audio.paused && !this.audio.ended && this.audio.readyState > 2)
  }

}
