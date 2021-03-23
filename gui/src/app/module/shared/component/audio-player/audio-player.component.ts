import { AfterViewInit, Component, ElementRef, Input, ViewChild } from '@angular/core';

@Component({
  selector: 'app-audio-player',
  templateUrl: './audio-player.component.html',
  styleUrls: ['./audio-player.component.scss']
})
export class AudioPlayerComponent implements AfterViewInit {

  @Input() public src: string;

  @Input() public autoplay: boolean = false;

  @Input() public showStateLabel: boolean = false;

  public audioStateLabel = 'Audio sample';

  @Input() public volume: number = 1.0; /* 1.0 is loudest */

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
    if (this.audio) {
      this.audio.volume = this.volume;
      this.audio.autoplay = this.autoplay;
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
    if (this.audio && this.audio.readyState >= 2) {
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

  public seek(time: number) {
    if (this.audio) {
      this.audio.currentTime = time;
    }
  }

  public playing(): boolean {
    return !!(this.audio.currentTime > 0 && !this.audio.paused && !this.audio.ended && this.audio.readyState > 2)
  }

}
