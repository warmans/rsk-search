import {BrowserModule, Title} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {RootRoutingModule} from './root-routing.module';
import {RootComponent} from './component/root/root.component';
import {CoreModule} from '../module/core/core.module';
import {HttpClientModule} from '@angular/common/http';
import {SearchAPIClientModule} from '../lib/api-client/services/search';
import {SearchModule} from '../module/search/search.module';
import {SharedModule} from '../module/shared/shared.module';
import {RewardModule} from '../module/reward/reward.module';
import {ContribModule} from '../module/contrib/contrib.module';
import {AdminModule} from '../module/admin/admin.module';
import {NgbModule} from '@ng-bootstrap/ng-bootstrap';
import {MoreShiteModule} from "../module/more-shite/more-shite.module";

@NgModule({
  declarations: [
    RootComponent
  ],
  imports: [
    CoreModule,
    BrowserModule,
    RootRoutingModule,
    HttpClientModule,
    SearchAPIClientModule.forRoot(),
    SharedModule,
    SearchModule,
    RewardModule,
    ContribModule,
    AdminModule,
    MoreShiteModule,
    NgbModule,
  ],
  providers: [Title],
  bootstrap: [RootComponent]
})
export class RootModule {
}
