import { enableProdMode } from '@angular/core';
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';
import { RootModule } from './app/root/root.module';
import { environment } from './environments/environment';
import * as Sentry from '@sentry/angular-ivy';

Sentry.init({
  dsn: 'https://8992b69abcda4231821c0697176ce365@o1428053.ingest.sentry.io/6777807',
  integrations: [
    new Sentry.BrowserTracing({
      tracingOrigins: ['localhost', 'https://scrimpton.com/api', 'https://karltakesasneakylookatmenscocks.com/api'],
      routingInstrumentation: Sentry.routingInstrumentation,
    }),
  ],

  // Set tracesSampleRate to 1.0 to capture 100%
  // of transactions for performance monitoring.
  // We recommend adjusting this value in production
  tracesSampleRate: 0.25,
});

if (environment.production) {
  enableProdMode();
}

platformBrowserDynamic().bootstrapModule(RootModule)
  .catch(err => console.error(err));
