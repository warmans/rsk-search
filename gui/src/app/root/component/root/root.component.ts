import {Component, EventEmitter, OnDestroy, OnInit, Renderer2} from '@angular/core';
import {Router} from '@angular/router';
import {Claims, SessionService} from 'src/app/module/core/service/session/session.service';
import {takeUntil} from 'rxjs/operators';
import {RskQuotas} from 'src/app/lib/api-client/models';
import {QuotaService} from 'src/app/module/core/service/quota/quota.service';
import {RadioService} from "../../../module/core/service/radio/radio.service";

@Component({
  selector: 'app-root',
  templateUrl: './root.component.html',
  styleUrls: ['./root.component.scss']
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
      this.bandwidthQuotaUsedPcnt = (1 - (res.bandwidthRemainingMib / res.bandwidthTotalMib));
    });
  }

  executeSearch(query: string) {
    this.router.navigate(['/search'], {queryParams: {q: query}});
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
