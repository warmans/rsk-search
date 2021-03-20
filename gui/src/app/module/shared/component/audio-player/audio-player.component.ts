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
    console.log(this.audioPlayerEl);
    if (!this.audioPlayerEl?.nativeElement) {
      return
    }
    this.audio = this.audioPlayerEl.nativeElement;
    if (this.audio) {
      this.audio.volume = this.volume;
      this.audio.autoplay = this.autoplay;
    }
  }

  public pause(): void {
    if (this.audio) {
      this.audio.pause();
      this.audioStateLabel = 'Paused';
    }
  }

  public get paused(): boolean {
    if (this.audio) {
      return this.audio.paused;
    } else {
      return true;
    }
  }

  public play(): void {
    console.log(this.audio);
    if (this.audio) {
      if (this.audio.readyState >= 2) {
        this.audio.play();
        this.audioStateLabel = 'Playing...';
      }
    }
  }


}
