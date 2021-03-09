/* tslint:disable */

import { HttpClient } from '@angular/common/http';
import { Inject, Injectable, Optional } from '@angular/core';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { DefaultHttpOptions, HttpOptions } from './';
import { USE_DOMAIN, USE_HTTP_OPTIONS, SearchAPIClient } from './search-api-client.service';

import * as models from '../../models';
import * as guards from '../../guards';

/**
 * Created with https://github.com/flowup/api-client-generator
 */
@Injectable()
export class GuardedSearchAPIClient extends SearchAPIClient {

  constructor(readonly httpClient: HttpClient,
              @Optional() @Inject(USE_DOMAIN) domain?: string,
              @Optional() @Inject(USE_HTTP_OPTIONS) options?: DefaultHttpOptions) {
    super(httpClient, domain, options);
  }

  searchServiceGetEpisode(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchEpisode> {
    return super.searchServiceGetEpisode(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchEpisode(res) || console.error(`TypeGuard for response 'RsksearchEpisode' caught inconsistency.`, res)));
  }

  searchServiceListFieldValues(
    args: {
      field?: string,
      prefix?: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchFieldValueList> {
    return super.searchServiceListFieldValues(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchFieldValueList(res) || console.error(`TypeGuard for response 'RsksearchFieldValueList' caught inconsistency.`, res)));
  }

  searchServiceSearch(
    args: {
      query?: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchResultList> {
    return super.searchServiceSearch(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskSearchResultList(res) || console.error(`TypeGuard for response 'RskSearchResultList' caught inconsistency.`, res)));
  }

}
