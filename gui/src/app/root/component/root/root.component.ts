import { Component, EventEmitter, OnDestroy, OnInit, Renderer2 } from '@angular/core';
import { Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { Claims, SessionService } from 'module/core/service/session/session.service';
import { catchError, finalize, map, switchMap, takeUntil } from 'rxjs/operators';
import { EMPTY } from 'rxjs';
import { RskQuotas, RskRandomQuote, RskTranscript } from 'lib/api-client/models';
import { SearchAPIClient } from 'lib/api-client/services/search';
import { QuotaService } from 'module/core/service/quota/quota.service';
import { AlertService } from 'module/core/service/alert/alert.service';
import { AudioService, PlayerMode } from 'module/core/service/audio/audio.service';
import { RadioService } from '../../../module/core/service/radio/radio.service';
import { SearchBarCompatComponent } from '../../../module/search/component/search-bar-compat/search-bar-compat.component';
import { NgClass, DecimalPipe } from '@angular/common';
import { UserMenuComponent } from 'module/shared/component/user-menu/user-menu.component';
import { AlertComponent } from 'module/shared/component/alert/alert.component';
import { AudioPlayerFixedComponent } from 'module/shared/component/audio-player-fixed/audio-player-fixed.component';
import { PendingRewardsComponent } from 'module/reward/component/pending-rewards/pending-rewards.component';

@Component({
  selector: 'app-root',
  templateUrl: './root.component.html',
  styleUrls: ['./root.component.scss'],
  imports: [
    RouterLink,
    SearchBarCompatComponent,
    NgClass,
    RouterLinkActive,
    RouterOutlet,
    DecimalPipe,
    UserMenuComponent,
    AlertComponent,
    AudioPlayerFixedComponent,
    PendingRewardsComponent,
  ],
})
export class RootComponent implements OnInit, OnDestroy {
  loggedInUser: Claims;

  darkTheme: boolean = true;

  destroy$: EventEmitter<boolean> = new EventEmitter<boolean>();

  quotas: RskQuotas;
  bandwidthQuotaUsedPcnt: number = 0;
  randomClipLoading: boolean = false;

  constructor(
    private renderer: Renderer2,
    private router: Router,
    private session: SessionService,
    private quotaService: QuotaService,
    private apiClient: SearchAPIClient,
    private audioService: AudioService,
    private alertService: AlertService,
    public radioService: RadioService,
  ) {
    session.onTokenChange.pipe(takeUntil(this.destroy$)).subscribe((token: string) => {
      if (token) {
        this.loggedInUser = this.session.getClaims();
      } else {
        this.loggedInUser = undefined;
      }
    });
    quotaService.quotas$.pipe(takeUntil(this.destroy$)).subscribe((res: RskQuotas) => {
      this.quotas = res;
      this.bandwidthQuotaUsedPcnt = 1 - res.bandwidthRemainingMib / res.bandwidthTotalMib;
    });
  }

  executeSearch(query: string) {
    this.router.navigate(['/search'], { queryParams: { q: query } });
  }

  logout() {
    this.session.destroySession();
    this.loggedInUser = undefined;
    this.router.navigate(['/search']);
  }

  ngOnDestroy(): void {
    this.destroy$.next(true);
    this.destroy$.complete();
  }

  toggleDarkmode() {
    this.darkTheme = !this.darkTheme;
    this.updateTheme();
  }

  updateTheme() {
    if (this.darkTheme) {
      this.renderer.removeClass(document.body, 'light-theme');
      this.renderer.addClass(document.body, 'dark-theme');
      localStorage.setItem('theme', 'dark');
    } else {
      this.renderer.removeClass(document.body, 'dark-theme');
      this.renderer.addClass(document.body, 'light-theme');
      localStorage.setItem('theme', 'light');
    }
  }

  ngOnInit(): void {
    this.darkTheme = (localStorage.getItem('theme') || 'dark') === 'dark';
    this.updateTheme();
  }

  toggleRadio() {
    if (!this.radioService.active) {
      this.radioService.start();
    } else {
      this.radioService.stop();
    }
  }

  playRandomClip() {
    if (this.randomClipLoading) {
      return;
    }

    this.randomClipLoading = true;
    this.apiClient
      .getRandomQuote()
      .pipe(
        switchMap((quote: RskRandomQuote) => {
          if (!quote.epid || quote.pos === undefined) {
            throw new Error('Random quote response did not include an episode and position.');
          }
          return this.apiClient.getTranscript({ epid: quote.epid }).pipe(map((transcript: RskTranscript) => ({ quote, transcript })));
        }),
        catchError((err) => {
          console.error('failed to play random clip', err);
          this.alertService.danger('Failed to play a random clip');
          return EMPTY;
        }),
        finalize(() => {
          this.randomClipLoading = false;
        }),
        takeUntil(this.destroy$),
      )
      .subscribe(({ quote, transcript }) => {
        const position = quote.pos;
        if (position === undefined) {
          return;
        }

        const line = transcript.transcript?.find((dialog) => dialog.pos === position) ?? transcript.transcript?.[position - 1];
        const offsetMs = Math.max(0, line?.offsetMs ?? 0);
        const episodeId = transcript.shortId || quote.epid;

        this.audioService.setAudioSrcFromEpisodeName(episodeId, transcript.name || episodeId, PlayerMode.Default);
        this.audioService.seekAudio(offsetMs / 1000);
        this.audioService.playAudio();
        this.router.navigate(['/ep', quote.epid], { fragment: `pos-${position}` });
      });
  }
}
