import { Component, EventEmitter, OnDestroy, OnInit, Renderer2 } from '@angular/core';
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router';
import { Claims, SessionService } from 'src/app/module/core/service/session/session.service';
import { takeUntil } from 'rxjs/operators';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { RskQuotas } from 'src/app/lib/api-client/models';
import { QuotaService } from 'src/app/module/core/service/quota/quota.service';

@Component({
  selector: 'app-root',
  templateUrl: './root.component.html',
  styleUrls: ['./root.component.scss']
})
export class RootComponent implements OnInit, OnDestroy {

  embedMode: boolean;

  loggedInUser: Claims;

  darkTheme: boolean = true;

  destory$: EventEmitter<boolean> = new EventEmitter<boolean>();

  quotas: RskQuotas;
  bandwidthQuotaUsedPcnt: number = 0;

  constructor(
    private renderer: Renderer2,
    private router: Router,
    private route: ActivatedRoute,
    private session: SessionService,
    private apiClient: SearchAPIClient,
    private quotaService: QuotaService,
  ) {
    session.onTokenChange.pipe(takeUntil(this.destory$)).subscribe((token) => {
      if (token) {
        this.loggedInUser = this.session.getClaims();
      } else {
        this.loggedInUser = undefined;
      }
    });
    this.router.events.pipe(takeUntil(this.destory$)).subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.embedMode = event.url.startsWith('/embed');
      }
    });
    quotaService.quotas$.pipe(takeUntil(this.destory$)).subscribe((res: RskQuotas) => {
      this.quotas = res;
      this.bandwidthQuotaUsedPcnt = (1 - (res.bandwidthRemainingMib / res.bandwidthTotalMib));
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
    this.destory$.next(true);
    this.destory$.complete();
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
    this.darkTheme = localStorage.getItem('theme') === 'dark';
    this.updateTheme();
  }
}
