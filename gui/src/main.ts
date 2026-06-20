import { enableProdMode, provideZoneChangeDetection, ErrorHandler, importProvidersFrom } from '@angular/core';
import { environment } from './environments/environment';
import * as Sentry from '@sentry/angular';
import { Title, bootstrapApplication } from '@angular/platform-browser';
import { provideHttpClient, withInterceptorsFromDi, HTTP_INTERCEPTORS } from '@angular/common/http';
import { provideRouter, withInMemoryScrolling, withRouterConfig } from '@angular/router';
import { routes } from './app/root/root.routes';
import { APIErrorInterceptor } from './app/module/core/interceptor/api-error.interceptor';
import { OutgoingTokenInterceptor } from './app/module/core/interceptor/outgoing-token.interceptor';
import { StarRatingModule } from 'angular-star-rating';
import { SearchAPIClientModule } from './app/lib/api-client/services/search';
import { CommunityAPIClientModule } from './app/lib/api-client/services/community';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { RootComponent } from './app/root/component/root/root.component';

Sentry.init({
  dsn: 'https://8992b69abcda4231821c0697176ce365@o1428053.ingest.sentry.io/6777807',
  integrations: [],

  // Set tracesSampleRate to 1.0 to capture 100%
  // of transactions for performance monitoring.
  // We recommend adjusting this value in production
  tracesSampleRate: 0.25,
});

if (environment.production) {
  enableProdMode();
}

bootstrapApplication(RootComponent, {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes, withInMemoryScrolling({ anchorScrolling: 'enabled' }), withRouterConfig({ onSameUrlNavigation: 'reload' })),
    importProvidersFrom(SearchAPIClientModule.forRoot(), CommunityAPIClientModule.forRoot(), StarRatingModule.forRoot(), NgbModule),
    Title,
    provideHttpClient(withInterceptorsFromDi()),
    { provide: HTTP_INTERCEPTORS, useClass: APIErrorInterceptor, multi: true },
    { provide: HTTP_INTERCEPTORS, useClass: OutgoingTokenInterceptor, multi: true },
    {
      provide: ErrorHandler,
      useValue: Sentry.createErrorHandler(),
    },
  ],
}).catch((err) => console.error(err));
