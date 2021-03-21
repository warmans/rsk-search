import { NgModule, Optional, SkipSelf } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AlertService } from './service/alert/alert.service';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { APIErrorInterceptor } from './interceptor/api-error.interceptor';
import { MetaService } from './service/meta/meta.service';
import { OutgoingTokenInterceptor } from './interceptor/outgoing-token.interceptor';

@NgModule({
  declarations: [],
  imports: [
    CommonModule
  ],
  providers: [
    AlertService,
    MetaService,
    { provide: HTTP_INTERCEPTORS, useClass: APIErrorInterceptor, multi: true },
    { provide: HTTP_INTERCEPTORS, useClass: OutgoingTokenInterceptor, multi: true },
  ]
})
export class CoreModule {

  constructor(@Optional() @SkipSelf() parentModule: CoreModule) {
    if (parentModule) {
      throw new Error(
        'CoreModule is already loaded. Import it in the AppModule only'
      );
    }
  }
}
