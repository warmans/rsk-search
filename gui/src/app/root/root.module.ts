import { BrowserModule, Title } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { RootRoutingModule } from './root-routing.module';
import { RootComponent } from './component/root/root.component';
import { CoreModule } from '../module/core/core.module';
import { HttpClientModule } from '@angular/common/http';
import { SearchAPIClientModule } from '../lib/api-client/services/search';
import { SearchModule } from '../module/search/search.module';
import { SharedModule } from '../module/shared/shared.module';
import { TscriptModule } from '../module/tscript/tscript.module';
import { RewardModule } from '../module/reward/reward.module';

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
    TscriptModule,
    RewardModule,
  ],
  providers: [Title],
  bootstrap: [RootComponent]
})
export class RootModule {
}
