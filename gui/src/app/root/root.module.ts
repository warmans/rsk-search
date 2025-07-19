import {BrowserModule, Title} from '@angular/platform-browser';
import {ErrorHandler, NgModule} from '@angular/core';

import {RootRoutingModule} from './root-routing.module';
import {RootComponent} from './component/root/root.component';
import {CoreModule} from '../module/core/core.module';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import {SearchAPIClientModule} from '../lib/api-client/services/search';
import {SearchModule} from '../module/search/search.module';
import {SharedModule} from '../module/shared/shared.module';
import {RewardModule} from '../module/reward/reward.module';
import {ContribModule} from '../module/contrib/contrib.module';
import {AdminModule} from '../module/admin/admin.module';
import {NgbModule} from '@ng-bootstrap/ng-bootstrap';
import {MoreShiteModule} from "../module/more-shite/more-shite.module";
import {CommunityAPIClientModule} from "../lib/api-client/services/community";

import * as Sentry from "@sentry/angular";

@NgModule({ declarations: [
        RootComponent
    ],
    bootstrap: [RootComponent],
    imports: [CoreModule,
        BrowserModule,
        RootRoutingModule,
        SearchAPIClientModule.forRoot(),
        CommunityAPIClientModule.forRoot(),
        SharedModule,
        SearchModule,
        RewardModule,
        ContribModule,
        AdminModule,
        MoreShiteModule,
        NgbModule],
    providers: [
      Title,
      provideHttpClient(withInterceptorsFromDi()),
      {
        provide: ErrorHandler,
        useValue: Sentry.createErrorHandler(),
      },
    ]
})
export class RootModule {
}
