import { Component, EventEmitter, OnDestroy, OnInit, Renderer2 } from '@angular/core';
import { Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { Claims, SessionService } from 'module/core/service/session/session.service';
import { takeUntil } from 'rxjs/operators';
import { RskQuotas } from 'lib/api-client/models';
import { QuotaService } from 'module/core/service/quota/quota.service';
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

  constructor(
    private renderer: Renderer2,
    private router: Router,
    private session: SessionService,
    private quotaService: QuotaService,
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
}
