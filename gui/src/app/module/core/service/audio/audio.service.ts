import { Injectable } from '@angular/core';
import { BehaviorSubject, combineLatest, Observable } from 'rxjs';

export interface Status {
  audioName: string;
  audioFile: string;
  state: PlayerState;
  currentTime: number;
  totalTime: number;
  percentElapsed: number;
  percentLoaded: number;
}

export interface TimeStatus {
  currentTime: number;
  totalTime: number;
  percentElapsed: number;
}

export interface FileStatus {
  audioName: string;
  audioFile: string;
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

  private audioName: string | null = null;

  private statusSub: BehaviorSubject<Status> = new BehaviorSubject<Status>({
    audioName: '',
    audioFile: '',
    state: PlayerState.paused,
    currentTime: 0,
    totalTime: 0,
    percentElapsed: 0,
    percentLoaded: 0,
  });
  public status: Observable<Status> = this.statusSub.asObservable();

  private timeStatusSub: BehaviorSubject<TimeStatus> = new BehaviorSubject<TimeStatus>({
    currentTime: 0,
    totalTime: 0,
    percentElapsed: 0
  });

  private playerStatusSub: BehaviorSubject<PlayerState> = new BehaviorSubject(PlayerState.paused);
  private percentLoadedSub: BehaviorSubject<number> = new BehaviorSubject(0);
  private audioSourceSub: BehaviorSubject<FileStatus> = new BehaviorSubject<FileStatus>({ audioFile: '', audioName: '' });

  constructor() {
    this.audio = new Audio();
    this.attachListeners();

    combineLatest([this.timeStatusSub, this.playerStatusSub, this.percentLoadedSub, this.audioSourceSub]).subscribe(([timeState, playerState, pcntLoaded, file]) => {

      if (file.audioFile == '') {
        this.statusSub.next(null);
        return;
      }

      const status = {
        audioName: file.audioName,
        audioFile: file.audioFile,
        state: playerState,
        currentTime: timeState.currentTime,
        totalTime: timeState.totalTime,
        percentElapsed: timeState.percentElapsed,
        percentLoaded: pcntLoaded,
      };
      this.statusSub.next(status);

      if (playerState === PlayerState.playing || playerState === PlayerState.paused || playerState === PlayerState.ended) {
        this.persistPlayerState(status);
      }
    });

    this.tryLoadPlayerState();
    this.tryLoadPlayerVolume();
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

  public reset() {
    this.pauseAudio();
    this.clearPlayerState();
    this.setAudioSrc(null, '');
  }

  /**
   * This is typically a URL to an MP3 file
   * @param name
   * @param src
   */
  public setAudioSrc(name: string | null, src: string): void {
    if (this.audioName === name) {
      return;
    }
    this.audioName = name;
    this.audio.src = src;
    this.audioSourceSub.next({ audioFile: src, audioName: name });
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

  public setVolume(vol: number): void {
    this.persistPlayerVolume(vol);
    this.audio.volume = vol;
  }

  private calculateTime = (evt) => {
    const ct = this.audio.currentTime;
    const d = this.audio.duration;

    this.timeStatusSub.next({
      currentTime: ct,
      totalTime: d,
      //todo: remove these
      percentElapsed: ((Math.floor((100 / d) * ct)) || 0),
    });
  };

  /**
   * This method takes the track's "duration" and "currentTime" properties to calculate the percent of time elapsed.
   * This is valuable for setting the position of a range input. It is called from the calculatePercentLoaded method.
   * @param p
   */
  private setPercentLoaded(p): void {
    this.percentLoadedSub.next(parseInt(p, 10) || 0);
  }

  private clearPlayerState() {
    localStorage.removeItem('audio_service_status');
  }

  private persistPlayerState(s: Status) {
    if (!s.audioFile || !s.audioName) {
      return;
    }
    localStorage.setItem('audio_service_status', JSON.stringify(s));
  }

  private tryLoadPlayerState() {
    const storedJSON = localStorage.getItem('audio_service_status');
    if (storedJSON) {
      let state: Status;
      try {
        state = JSON.parse(storedJSON);
      } catch (e) {
        console.error(`failed to load audio service state: ${e}`);
      }
      if (state) {
        this.setAudioSrc(state.audioName, state.audioFile);
        this.audio.load();
        this.seekAudio(state.currentTime);
      }
    }
  }

  private persistPlayerVolume(vol: number) {
    localStorage.setItem('audio_service_volume', JSON.stringify(vol));
  }

  private tryLoadPlayerVolume() {
    const storedJSON = localStorage.getItem('audio_service_volume');
    let vol: number = 1;
    if (storedJSON) {
      try {
        vol = JSON.parse(storedJSON);
      } catch (e) {
        console.error(`failed to load audio service state: ${e}`);
      }
    }
    this.audio.volume = vol || 1;
  }
}
