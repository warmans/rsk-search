import {Injectable, Query} from '@angular/core';
import {BehaviorSubject, combineLatest, Observable} from 'rxjs';
import {HttpParams} from "@angular/common/http";

const STORAGE_KEY_LISTENLOG = 'audio_service_listen_log';
const STORAGE_KEY_VOLUME = 'audio_service_volume';

export enum PlayerMode {
  Default = 'default',
  Standalone = 'standalone',
  Radio = 'radio'
}

export interface Status {
  audioID: string;
  audioName: string,
  audioFile: string;
  mode: PlayerMode;
  state: PlayerState;
  currentTime: number;
  totalTime: number;
  percentElapsed: number;
  percentLoaded: number;
  listened: boolean;
  sleepTimerRemainder: number | undefined;
}

export interface TimeStatus {
  currentTime: number;
  totalTime: number;
  percentElapsed: number;
}

export interface FileStatus {
  audioID: string;
  audioName: string,
  audioFile: string;
  mode: PlayerMode;
}

export enum PlayerState {
  playing = 'playing',
  paused = 'paused',
  loading = 'loading',
  ended = 'ended',
  failed = 'failed',
}

@Injectable({
  providedIn: 'root'
})
export class AudioService {

  public audio: HTMLAudioElement;

  private audioID: string | null = null;

  private statusSub: BehaviorSubject<Status> = new BehaviorSubject<Status>({
    audioID: '',
    audioName: '',
    audioFile: '',
    mode: PlayerMode.Default,
    state: PlayerState.paused,
    currentTime: 0,
    totalTime: 0,
    percentElapsed: 0,
    percentLoaded: 0,
    listened: false,
    sleepTimerRemainder: undefined,
  });
  public status: Observable<Status> = this.statusSub.asObservable();

  private volumeLoadedSub: BehaviorSubject<number> = new BehaviorSubject<number>(1);
  public volumeLoaded: Observable<number> = this.volumeLoadedSub.asObservable();

  private modeSub: BehaviorSubject<PlayerMode> = new BehaviorSubject<PlayerMode>(PlayerMode.Default);
  public mode$: Observable<PlayerMode> = this.modeSub.asObservable();

  private timeStatusSub: BehaviorSubject<TimeStatus> = new BehaviorSubject<TimeStatus>({
    currentTime: 0,
    totalTime: 0,
    percentElapsed: 0
  });

  private playerStatusSub: BehaviorSubject<PlayerState> = new BehaviorSubject(PlayerState.paused);
  private percentLoadedSub: BehaviorSubject<number> = new BehaviorSubject(0);
  private audioSourceSub: BehaviorSubject<FileStatus> = new BehaviorSubject<FileStatus>({
    audioFile: '',
    audioID: '',
    audioName: '',
    mode: PlayerMode.Default
  });
  private audioHistoryLogSub: BehaviorSubject<string[]> = new BehaviorSubject<string[]>(this.getListenLog());
  private errorsSub: BehaviorSubject<string[]> = new BehaviorSubject<string[]>(this.getListenLog());

  public audioHistoryLog = this.audioHistoryLogSub.asObservable();

  private sleepTimerRemaining: number = undefined;
  private sleepInterval: any;

  constructor() {
    this.audio = new Audio();
    this.attachListeners();

    this.playerStatusSub.subscribe((playerState) => {
      if (playerState === PlayerState.ended) {
        this.persistEpisodeListened(this.audioID);
        if (this.modeSub.getValue() !== 'radio') {
          this.clearSleepTimer();
        }
      }
    })

    combineLatest(
      [this.timeStatusSub, this.playerStatusSub, this.percentLoadedSub, this.audioSourceSub, this.audioHistoryLogSub]
    ).subscribe(([timeState, playerState, pcntLoaded, file, history]) => {
      const status = {
        audioID: file.audioID,
        audioName: file.audioName,
        audioFile: file.audioFile,
        mode: file.mode,
        state: playerState,
        currentTime: timeState.currentTime,
        totalTime: timeState.totalTime,
        percentElapsed: timeState.percentElapsed,
        percentLoaded: pcntLoaded,
        listened: this.modeSub.getValue() === PlayerMode.Default ? history.indexOf(file.audioID) > -1 : false,
        volume: this.audio.volume,
        sleepTimerRemainder: this.sleepTimerRemaining,
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
    this.audio.addEventListener('error', this.setPlayerStatus, false);
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
      case 'error':
        this.playerStatusSub.next(PlayerState.failed);
        break;
      default:
        this.playerStatusSub.next(PlayerState.paused);
        break;
    }
  };

  public reset() {
    this.pauseAudio();
    this.setAudioSrc(null, '', this.modeSub.getValue());
    this.clearPersistentPlayerState();
  }

  public setAudioSrc(id: string | null, name: string | null, mode?: PlayerMode, startMs?: number, endMs?: number): void {
    if (this.audioID === id) {
      return;
    }
    this.pauseAudio();

    let query = new HttpParams()
    if (startMs || endMs) {
      query = query.set("ts", `${startMs}${endMs ? "-" + endMs : ""}`)
    }
    if (mode == 'radio') {
      query = query.set("remastered", "1")
    }
    let audioUri: string = `/dl/media/${id}.mp3?` + query.toString();

    this.audioID = id;
    if (id !== null) {
      this.audio.src = audioUri
    }
    this.modeSub.next(mode ?? PlayerMode.Default);
    this.audioSourceSub.next({audioFile: audioUri, audioID: id, audioName: name, mode: mode});
  }

  public playAudio(withOffset?: number): void {
    if (!this.audio) {
      return;
    }
    if (withOffset !== undefined) {
      this.audio.currentTime = this.audio.currentTime + withOffset > 0 ? this.audio.currentTime + withOffset : 0;
    }
    // if this is a regular player always use the full playback rate
    if (this.modeSub.getValue() !== PlayerMode.Standalone) {
      this.setPlaybackRate(1);
    }
    this.audio.play();
  }

  public pauseAudio(): void {
    // force an event to be emitted ASAP.
    this.playerStatusSub.next(PlayerState.paused);
    this.audio.pause();
  }

  public toggleAudio(withOffset?: number): void {
    (this.audio.paused) ? this.playAudio(withOffset) : this.pauseAudio();
  }

  /**
   * @param position number seconds.milliseconds
   */
  public seekAudio(position: number): void {
    this.audio.currentTime = position;
  }

  public setVolume(vol: number): void {
    this.persistPlayerVolume(vol);
    this.audio.volume = vol;
  }

  public setPlaybackRate(rate: number): void {
    this.audio.defaultPlaybackRate = rate;
    this.audio.playbackRate = rate;
  }

  public getListenLog(): string[] {
    let listenLog: string[] = [];
    try {
      listenLog = JSON.parse(localStorage.getItem(STORAGE_KEY_LISTENLOG) || '[]');
    } catch (e) {
    }
    return listenLog;
  }

  public getCurrentFileIsListened(): boolean {
    return this.getListenLog().indexOf(this.audioID) > -1;
  }

  public markAsPlayed(): void {
    this.persistEpisodeListened(this.audioID);
    this.playerStatusSub.next(PlayerState.ended);
    if (this.modeSub.getValue() === PlayerMode.Default){
      this.reset();
    }
  }

  public markAsUnplayed(): void {
    this.persistEpisodeUnlistened(this.audioID);
  }

  public incrementSleepTimer(durationSeconds: number) {
    if (this.sleepTimerRemaining == null) {
      this.sleepTimerRemaining = 0;
    }
    this.sleepTimerRemaining += 1000 * durationSeconds;
    if (this.sleepInterval) {
      return;
    }
    this.sleepInterval = setInterval(() => {
      if (this.sleepTimerRemaining <= 0) {
        this.clearSleepTimer();
        this.pauseAudio();
      } else {
        this.sleepTimerRemaining -= 1000;
      }
    }, 1000)
  }

  public clearSleepTimer() {
    this.sleepTimerRemaining = 0;
    clearInterval(this.sleepInterval);
    this.sleepInterval = undefined;
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

  private setPercentLoaded(p): void {
    this.percentLoadedSub.next(parseInt(p, 10) || 0);
  }

  private clearPersistentPlayerState() {
    if (this.modeSub.getValue() === PlayerMode.Default) {
      localStorage.removeItem(this.statusStorageKey());
    }
  }

  private persistPlayerState(s: Status) {
    if (!s.audioFile || !s.audioID || this.modeSub.getValue() !== PlayerMode.Default) {
      return;
    }
    localStorage.setItem(this.statusStorageKey(), JSON.stringify(s));
  }

  private persistEpisodeListened(audioID: string) {
    if (this.modeSub.getValue() !== PlayerMode.Default) {
      return;
    }

    const listenLog = this.getListenLog()
    if (listenLog.indexOf(audioID) === -1) {
      listenLog.push(audioID);
    }
    localStorage.setItem(STORAGE_KEY_LISTENLOG, JSON.stringify(listenLog));
    this.audioHistoryLogSub.next(listenLog);
  }

  private persistEpisodeUnlistened(audioID: string) {
    if (this.modeSub.getValue() !== PlayerMode.Default) {
      return;
    }
    const listenLog = this.getListenLog();
    const nameIdx = listenLog.indexOf(audioID);
    if (nameIdx > -1) {
      listenLog.splice(nameIdx, 1);
      localStorage.setItem(STORAGE_KEY_LISTENLOG, JSON.stringify(listenLog));
      this.audioHistoryLogSub.next(listenLog);
    }
  }

  private tryLoadPlayerState() {
    const storedJSON = localStorage.getItem(this.statusStorageKey());
    if (storedJSON) {
      let state: Status;
      try {
        state = JSON.parse(storedJSON);
      } catch (e) {
        console.error(`failed to load audio service state: ${e}`);
      }

      if (state && state.mode === PlayerMode.Default) {
        this.setAudioSrc(state.audioID, state.audioName);
        this.audio.load();
        this.seekAudio(state.currentTime);
      }
    }
  }

  private persistPlayerVolume(vol: number) {
    localStorage.setItem(STORAGE_KEY_VOLUME, JSON.stringify(vol));
  }

  private tryLoadPlayerVolume() {
    const storedJSON = localStorage.getItem(STORAGE_KEY_VOLUME);
    let vol: number = 1;
    if (storedJSON) {
      try {
        vol = JSON.parse(storedJSON);
      } catch (e) {
        console.error(`failed to load audio service state: ${e}`);
      }
    }
    this.audio.volume = vol ?? 1;
    this.volumeLoadedSub.next(this.audio.volume);
  }

  private statusStorageKey() {
    return `audio_service_status${this.modeSub.getValue() === PlayerMode.Standalone ? '-temp' : ''}`
  }
}
