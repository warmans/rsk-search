import { Component, EventEmitter, Input, OnDestroy, OnInit } from '@angular/core';
import { takeUntil } from 'rxjs/operators';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { SessionService } from '../../../core/service/session/session.service';

@Component({
  selector: 'app-reddit-login',
  templateUrl: './reddit-login.component.html',
  styleUrls: ['./reddit-login.component.scss']
})
export class RedditLoginComponent implements OnInit, OnDestroy {

  @Input()
  open: boolean = false;

  showMoreAuthInformation: boolean = false;

  authError: string;

  loading: boolean;

  authenticated: boolean = false;

  destroy$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(
    private apiClient: SearchAPIClient,
    private route: ActivatedRoute,
    private router: Router,
    private sessionService: SessionService,
  ) {
    route.queryParamMap.pipe(takeUntil(this.destroy$)).subscribe((d: Data) => {
      this.authError = d.params['error'];

      if (d.params['token']) {
        this.sessionService.registerToken(d.params['token']);

        // once a token has been stored clear it from the URL
        let urlTree = this.router.parseUrl(this.router.url);
        urlTree.queryParams = {};
        urlTree.fragment = null;
        this.router.navigate([urlTree.toString()]);
      }
    });

    this.sessionService.onTokenChange.pipe(takeUntil(this.destroy$)).subscribe((token: string): void => {
      if (token != null) {
        this.authenticated = true;
      }
    });
  }

  ngOnInit(): void {

  }

  requestAuth() {
    this.loading = true;
    this.apiClient.searchServiceGetRedditAuthURL().pipe(takeUntil(this.destroy$)).subscribe((res) => {
      document.location.href = res.url;
    }).add(() => this.loading = false);
  }

  ngOnDestroy(): void {
    this.destroy$.next(true);
    this.destroy$.complete();
  }

}
