import { Component, EventEmitter, OnDestroy, OnInit, Renderer2 } from '@angular/core';
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router';
import { Claims, SessionService } from 'src/app/module/core/service/session/session.service';
import { takeUntil } from 'rxjs/operators';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { RskPrediction, RskSearchTermPredictions } from 'src/app/lib/api-client/models';

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

  searchPredictions: RskPrediction[] = [];

  constructor(
    private renderer: Renderer2,
    private router: Router,
    private route: ActivatedRoute,
    private session: SessionService,
    private apiClient: SearchAPIClient,
  ) {
    session.onTokenChange.pipe(takeUntil(this.destory$)).subscribe((token) => {
      if (token) {
        this.loggedInUser = this.session.getClaims();
      } else {
        this.loggedInUser = undefined;
      }
    });
    this.router.events.pipe(takeUntil(this.destory$)).subscribe((event) => {
      if(event instanceof NavigationEnd) {
        this.embedMode = event.url.startsWith("/embed")
      }
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

  predictTerms(prefix: string) {
    this.searchPredictions = [];
    if ((prefix || '').trim() == '') {
      return;
    }
    this.apiClient.predictSearchTerm({ prefix: prefix, maxPredictions: 5 })
      .pipe(takeUntil(this.destory$))
      .subscribe((value: RskSearchTermPredictions) => {
        this.searchPredictions = value.predictions;
      });
  }
}
