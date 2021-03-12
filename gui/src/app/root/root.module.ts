import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { RootRoutingModule } from './root-routing.module';
import { RootComponent } from './component/root/root.component';
import { CoreModule } from '../module/core/core.module';
import { HttpClientModule } from '@angular/common/http';
import { SearchAPIClientModule } from '../lib/api-client/services/search';
import { SearchModule } from '../module/search/search.module';
import { ReactiveFormsModule } from '@angular/forms';

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
    SearchModule,
  ],
  providers: [],
  bootstrap: [RootComponent]
})
export class RootModule {
}
