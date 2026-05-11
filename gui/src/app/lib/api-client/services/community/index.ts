/* tslint:disable */

import { NgModule, ModuleWithProviders } from '@angular/core';
import { HttpHeaders, HttpParams } from '@angular/common/http';
import { CommunityAPIClient, USE_DOMAIN, USE_HTTP_OPTIONS } from './community-api-client.service';
import { GuardedCommunityAPIClient } from './guarded-community-api-client.service';

export { CommunityAPIClient } from './community-api-client.service';
export { CommunityAPIClientInterface } from './community-api-client.interface';
export { GuardedCommunityAPIClient } from './guarded-community-api-client.service';

/**
 * provided options, headers and params will be used as default for each request
 */
export interface DefaultHttpOptions {
  headers?: {[key: string]: string};
  params?: {[key: string]: string};
  reportProgress?: boolean;
  withCredentials?: boolean;
}

export interface HttpOptions {
  headers?: HttpHeaders;
  params?: HttpParams;
  reportProgress?: boolean;
  withCredentials?: boolean;
}

export interface CommunityAPIClientModuleConfig {
  domain?: string;
  guardResponses?: boolean; // validate responses with type guards
  httpOptions?: DefaultHttpOptions;
}

@NgModule({})
export class CommunityAPIClientModule {
  /**
   * Use this method in your root module to provide the CommunityAPIClientModule
   *
   * If you are not providing
   * @param { CommunityAPIClientModuleConfig } config
   * @returns { ModuleWithProviders }
   */
  static forRoot(config: CommunityAPIClientModuleConfig = {}): ModuleWithProviders<CommunityAPIClientModule> {
    return {
      ngModule: CommunityAPIClientModule,
      providers: [
        ...(config.domain != null ? [{provide: USE_DOMAIN, useValue: config.domain}] : []),
        ...(config.httpOptions ? [{provide: USE_HTTP_OPTIONS, useValue: config.httpOptions}] : []),
        ...(config.guardResponses ? [{provide: CommunityAPIClient, useClass: GuardedCommunityAPIClient }] : [CommunityAPIClient]),
      ]
    };
  }
}
