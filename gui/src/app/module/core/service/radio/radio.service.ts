import {Injectable} from '@angular/core';
import {AudioService, PlayerMode, PlayerState} from "../audio/audio.service";
import {takeUntil} from "rxjs/operators";
import {Subject} from "rxjs";
import {SearchAPIClient} from "../../../../lib/api-client/services/search";
import {RskRadioState} from "../../../../lib/api-client/models";
import {SessionService} from "../session/session.service";

@Injectable({
  providedIn: 'root'
})
export class RadioService {

  private detatchAudioService: Subject<void> = new Subject<void>();

  private lastUpdateSent: number = 0;

  private userIsLoggedIn: boolean = false;

  private state: RskRadioState = null;

  public active: boolean = false;

  constructor(private audioService: AudioService, private apiClient: SearchAPIClient, private session: SessionService) {
    this.userIsLoggedIn = session.getToken() != null;

    this.session.onTokenChange.subscribe((token) => {
      this.userIsLoggedIn = token != null;
    })

    this.audioService.mode$.subscribe((newMode: PlayerMode) => {
      if (newMode === PlayerMode.Radio) {
        this.attach();
        return;
      }
      this.detatch();
    });
  }

  start() {
    if (!this.userIsLoggedIn) {
      return
    }
    this.apiClient.getState().subscribe((state) => {
      this.state = state;
      this.applyState();
      this.attach();
    })
  }

  attach() {
    if (this.active === true) {
      return;
    }

    this.active = true;
    this.audioService.status.pipe(takeUntil(this.detatchAudioService)).subscribe((status) => {
      // player reset or closed or something.
      if (status.audioID === null) {
        this.detatch();
        return;
      }
      // episode is finished
      if (status.state === PlayerState.ended) {
        this.fetchNext();
        return;
      }
      // something else happened
      if (status.state !== PlayerState.playing || !this.state) {
        return;
      }

      // episode is playing
      this.state.currentTimestampMs = Math.floor(status.currentTime);
      if (status.audioID !== this.state.currentEpisode.shortId) {
        this.state.currentEpisode.shortId = status.audioID;
        this.state.currentEpisode.startedAt = (new Date()).toISOString();
      }

      const nowTs = (new Date()).getTime();
      if (nowTs - this.lastUpdateSent > 2500 || this.lastUpdateSent === 0) {
        this.lastUpdateSent = nowTs;
        this.storeState();
      }
    });
  }

  detatch() {
    this.active = false;
    this.detatchAudioService.next();
  }

  stop() {
    this.audioService.reset();
  }

  applyState() {
    if (!this.state) {
      this.audioService.reset();
      return;
    }
    this.audioService.setAudioSrc(
      this.state.currentEpisode.shortId,
      "RADIO",
      PlayerMode.Radio,
    );

    this.lastUpdateSent = 0;
    this.audioService.seekAudio(this.state.currentTimestampMs ?? 0);
    this.audioService.playAudio()
  }

  fetchNext() {
    if (!this.userIsLoggedIn) {
      return;
    }
    this.apiClient.getNext().subscribe((next) => {

      this.state = {
        currentEpisode: {
          shortId: next.shortId,
          startedAt: (new Date()).toISOString(),
        },
        currentTimestampMs: 0,
      };

      this.applyState();
    })
  }

  storeState() {
    if (!this.userIsLoggedIn) {
      return;
    }
    if (!this.state) {
      return;
    }
    this.apiClient.putState({body: this.state}).subscribe();
  }
}
