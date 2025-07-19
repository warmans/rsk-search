import {Component, EventEmitter, Input, OnDestroy} from '@angular/core';
import {takeUntil} from 'rxjs/operators';
import {SearchAPIClient} from 'src/app/lib/api-client/services/search';
import {ActivatedRoute, Data, Router} from '@angular/router';
import {SessionService} from 'src/app/module/core/service/session/session.service';
import {FormControl} from "@angular/forms";
import {RskAuthURL} from "../../../../lib/api-client/models";

@Component({
    selector: 'app-oauth-login',
    templateUrl: './oauth-login.component.html',
    styleUrls: ['./oauth-login.component.scss'],
    standalone: false
})
export class OauthLoginComponent implements OnDestroy {

  @Input()
  open: boolean = false;

  authMethod: FormControl<string> = new FormControl<string>("reddit");

  showMoreAuthInformation: boolean = false;

  authError: string;

  loading: boolean;

  authenticated: boolean = false;

  destroy$: EventEmitter<void> = new EventEmitter<void>();

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

  requestAuth() {
    this.loading = true;
    this.apiClient.getAuthUrl({provider: this.authMethod.value}).pipe(takeUntil(this.destroy$)).subscribe((res: RskAuthURL) => {
      document.location.href = res.url;
    }).add((): boolean => this.loading = false);
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

}
