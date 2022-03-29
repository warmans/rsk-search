import { Injectable } from '@angular/core';
import { BehaviorSubject, combineLatest, Observable } from 'rxjs';

// via: https://gist.github.com/philmerrell/d65655ef73a5b3be863491b19b3902ba

export interface Status {
  audioName: string;
  state: PlayerState;
  duration: string;
  currentTime: number;
  timeElapsed: string;
  timeRemaining: string;
  percentElapsed: number;
  percentLoaded: number;
}

export interface TimeStatus {
  duration: string,
  currentTime: number,
  timeElapsed: string;
  timeRemaining: string;
  percentElapsed: number;
}

export enum PlayerState {
  playing = 'playing',
  paused = 'paused',
  loading = 'loading',
  ended = 'ended',
}

@Injectable({
  providedIn: 'root'
})
export class AudioService {

  public audio: HTMLAudioElement;

  private audioName: string = '';

  private statusSub: BehaviorSubject<Status> = new BehaviorSubject<Status>({
    audioName: '',
    state: PlayerState.paused,
    duration: '00:00',
    currentTime: 0,
    percentElapsed: 0,
    timeElapsed: '00:00',
    timeRemaining: '-00:00',
    percentLoaded: 0,
  });
  public status: Observable<Status> = this.statusSub.asObservable();

  private timeStatusSub: BehaviorSubject<TimeStatus> = new BehaviorSubject({ currentTime: 0, duration: '00:00', timeElapsed: '00:00', timeRemaining: '-00:00', percentElapsed: 0 });
  private playerStatusSub: BehaviorSubject<PlayerState> = new BehaviorSubject(PlayerState.paused);
  private percentLoadedSub: BehaviorSubject<number> = new BehaviorSubject(0);

  constructor() {
    this.audio = new Audio();
    this.attachListeners();

    combineLatest([this.timeStatusSub, this.playerStatusSub, this.percentLoadedSub]).subscribe(([timeState, playerState, pcntLoaded]) => {
      this.statusSub.next({
        audioName: this.audioName,
        state: playerState,
        duration: timeState.duration,
        currentTime: timeState.currentTime,
        percentElapsed: timeState.percentElapsed,
        timeElapsed: timeState.timeElapsed,
        timeRemaining: timeState.timeRemaining,
        percentLoaded: pcntLoaded,
      });
    });
  }

  private attachListeners(): void {
    this.audio.addEventListener('timeupdate', this.calculateTime, false);
    this.audio.addEventListener('playing', this.setPlayerStatus, false);
    this.audio.addEventListener('pause', this.setPlayerStatus, false);
    this.audio.addEventListener('progress', this.calculatePercentLoaded, false);
    this.audio.addEventListener('waiting', this.setPlayerStatus, false);
    this.audio.addEventListener('ended', this.setPlayerStatus, false);
  }

  private calculatePercentLoaded = (evt) => {
    if (this.audio.duration > 0) {
      for (var i = 0; i < this.audio.buffered.length; i++) {
        if (this.audio.buffered.start(this.audio.buffered.length - 1 - i) < this.audio.currentTime) {
          let percent = (this.audio.buffered.end(this.audio.buffered.length - 1 - i) / this.audio.duration) * 100;
          this.setPercentLoaded(percent);
          break;
        }
      }
    }
  };

  private setPlayerStatus = (evt) => {
    switch (evt.type) {
      case 'playing':
        this.playerStatusSub.next(PlayerState.playing);
        break;
      case 'pause':
        this.playerStatusSub.next(PlayerState.paused);
        break;
      case 'waiting':
        this.playerStatusSub.next(PlayerState.loading);
        break;
      case 'ended':
        this.playerStatusSub.next(PlayerState.ended);
        break;
      default:
        this.playerStatusSub.next(PlayerState.paused);
        break;
    }
  };

  /**
   * If you need the audio instance in your component for some reason, use this.
   */
  public getAudio(): HTMLAudioElement {
    return this.audio;
  }

  /**
   * This is typically a URL to an MP3 file
   * @param src
   */
  public setAudioSrc(name: string, src: string): void {
    this.audioName = name;
    this.audio.src = src;
  }

  /**
   * The method to play audio
   */
  public playAudio(): void {
    this.audio.play();
  }

  /**
   * The method to pause audio
   */
  public pauseAudio(): void {
    this.audio.pause();
  }

  /**
   * Convenience method to toggle the audio between playing and paused
   */
  public toggleAudio(): void {
    (this.audio.paused) ? this.audio.play() : this.audio.pause();
  }

  /**
   * Method to seek to a position on the audio track (in milliseconds, I think),
   * @param position
   */
  public seekAudio(position: number): void {
    this.audio.currentTime = position;
  }

  private calculateTime = (evt) => {
    const ct = this.audio.currentTime;
    const d = this.audio.duration;

    this.timeStatusSub.next({
      duration: this.formatInterval(d),
      currentTime: ct,
      timeElapsed: this.formatInterval(ct),
      timeRemaining: this.calculateTimeRemaining(d, ct),
      percentElapsed: ((Math.floor((100 / d) * ct)) || 0),
    });
  };

  /**
   * This formats the audio's elapsed time into a human readable format, should be refactored into a Pipe.
   * It takes the audio track's "currentTime" property as an argument. It is called from the, calulateTime method.
   * @param interval interval seconds.subseconds
   */
  private formatInterval(interval: number): string {
    let seconds = Math.floor(interval % 60),
      displaySecs = (seconds < 10) ? '0' + seconds : seconds,
      minutes = Math.floor((interval / 60) % 60),
      displayMins = (minutes < 10) ? '0' + minutes : minutes;

    return `${displayMins}:${displaySecs}`;
  }


  /**
   * This method takes the track's "duration" and "currentTime" properties to calculate the remaing time the track has
   * left to play.
   * @param d
   * @param t
   */
  private calculateTimeRemaining(d: number, t: number): string {
    let remaining;
    let timeLeft = d - t,
      seconds = Math.floor(timeLeft % 60) || 0,
      remainingSeconds = seconds < 10 ? '0' + seconds : seconds,
      minutes = Math.floor((timeLeft / 60) % 60) || 0,
      remainingMinutes = minutes < 10 ? '0' + minutes : minutes,
      hours = Math.floor(((timeLeft / 60) / 60) % 60) || 0;

    // remaining = (hours === 0)
    if (hours === 0) {
      remaining = '-' + remainingMinutes + ':' + remainingSeconds;
    } else {
      remaining = '-' + hours + ':' + remainingMinutes + ':' + remainingSeconds;
    }
    return remaining;
  }

  /**
   * This method takes the track's "duration" and "currentTime" properties to calculate the percent of time elapsed.
   * This is valuable for setting the position of a range input. It is called from the calculatePercentLoaded method.
   * @param p
   */
  private setPercentLoaded(p): void {
    this.percentLoadedSub.next(parseInt(p, 10) || 0);
  }

}
