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

  searchServiceGetRedditAuthURL(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchRedditAuthURL> {
    return super.searchServiceGetRedditAuthURL(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchRedditAuthURL(res) || console.error(`TypeGuard for response 'RsksearchRedditAuthURL' caught inconsistency.`, res)));
  }

  searchServiceAuthorizeRedditToken(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchToken> {
    return super.searchServiceAuthorizeRedditToken(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchToken(res) || console.error(`TypeGuard for response 'RsksearchToken' caught inconsistency.`, res)));
  }

  searchServiceListEpisodes(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchEpisodeList> {
    return super.searchServiceListEpisodes(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchEpisodeList(res) || console.error(`TypeGuard for response 'RsksearchEpisodeList' caught inconsistency.`, res)));
  }

  searchServiceSubmitDialogCorrection(
    args: {
      episodeId: string,
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    return super.searchServiceSubmitDialogCorrection(args, requestHttpOptions)
      .pipe(tap((res: any) => typeof res === 'object' || console.error(`TypeGuard for response 'object' caught inconsistency.`, res)));
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

  searchServiceGetSearchMetadata(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchMetadata> {
    return super.searchServiceGetSearchMetadata(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskSearchMetadata(res) || console.error(`TypeGuard for response 'RskSearchMetadata' caught inconsistency.`, res)));
  }

  searchServiceSearch(
    args: {
      query?: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchResultList> {
    return super.searchServiceSearch(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskSearchResultList(res) || console.error(`TypeGuard for response 'RskSearchResultList' caught inconsistency.`, res)));
  }

  searchServiceSubmitTscriptChunk(
    args: {
      chunkId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    return super.searchServiceSubmitTscriptChunk(args, requestHttpOptions)
      .pipe(tap((res: any) => typeof res === 'object' || console.error(`TypeGuard for response 'object' caught inconsistency.`, res)));
  }

  searchServiceListTscriptChunkSubmissions(
    args: {
      chunkId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkSubmissionList> {
    return super.searchServiceListTscriptChunkSubmissions(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchChunkSubmissionList(res) || console.error(`TypeGuard for response 'RsksearchChunkSubmissionList' caught inconsistency.`, res)));
  }

  searchServiceGetTscriptChunk(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptChunk> {
    return super.searchServiceGetTscriptChunk(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchTscriptChunk(res) || console.error(`TypeGuard for response 'RsksearchTscriptChunk' caught inconsistency.`, res)));
  }

  searchServiceGetTscriptChunkStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkStats> {
    return super.searchServiceGetTscriptChunkStats(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchChunkStats(res) || console.error(`TypeGuard for response 'RsksearchChunkStats' caught inconsistency.`, res)));
  }

  searchServiceListFieldValues(
    args: {
      field: string,
      prefix?: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchFieldValueList> {
    return super.searchServiceListFieldValues(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchFieldValueList(res) || console.error(`TypeGuard for response 'RsksearchFieldValueList' caught inconsistency.`, res)));
  }

}
