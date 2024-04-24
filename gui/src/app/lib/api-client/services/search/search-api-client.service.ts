/* tslint:disable */

import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Inject, Injectable, InjectionToken, Optional } from '@angular/core';
import { Observable, throwError } from 'rxjs';
import { DefaultHttpOptions, HttpOptions, SearchAPIClientInterface } from './';

import * as models from '../../models';

export const USE_DOMAIN = new InjectionToken<string>('SearchAPIClient_USE_DOMAIN');
export const USE_HTTP_OPTIONS = new InjectionToken<HttpOptions>('SearchAPIClient_USE_HTTP_OPTIONS');

type APIHttpOptions = HttpOptions & {
  headers: HttpHeaders;
  params: HttpParams;
  responseType?: 'arraybuffer' | 'blob' | 'text' | 'json';
};

/**
 * Created with https://github.com/flowup/api-client-generator
 */
@Injectable()
export class SearchAPIClient implements SearchAPIClientInterface {

  readonly options: APIHttpOptions;

  readonly domain: string = `//${window.location.hostname}${window.location.port ? ':'+window.location.port : ''}`;

  constructor(private readonly http: HttpClient,
              @Optional() @Inject(USE_DOMAIN) domain?: string,
              @Optional() @Inject(USE_HTTP_OPTIONS) options?: DefaultHttpOptions) {

    if (domain != null) {
      this.domain = domain;
    }

    this.options = {
      headers: new HttpHeaders(options && options.headers ? options.headers : {}),
      params: new HttpParams(options && options.params ? options.params : {}),
      ...(options && options.reportProgress ? { reportProgress: options.reportProgress } : {}),
      ...(options && options.withCredentials ? { withCredentials: options.withCredentials } : {})
    };
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  createTscriptImport(
    args: {
      body: models.RskCreateTscriptImportRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTscriptImport> {
    const path = `/api/admin/tscript/import`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskTscriptImport>('POST', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listTscriptImports(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTscriptImportList> {
    const path = `/api/admin/tscript/imports`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('filter' in args) {
      options.params = options.params.set('filter', String(args.filter));
    }
    if ('sortField' in args) {
      options.params = options.params.set('sortField', String(args.sortField));
    }
    if ('sortDirection' in args) {
      options.params = options.params.set('sortDirection', String(args.sortDirection));
    }
    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    if ('pageSize' in args) {
      options.params = options.params.set('pageSize', String(args.pageSize));
    }
    return this.sendRequest<models.RskTscriptImportList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  deleteTscript(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    const path = `/api/admin/tscript/${args.id}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<object>('DELETE', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getAuthUrl(
    args: {
      provider?: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskAuthURL> {
    const path = `/api/auth/url`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('provider' in args) {
      options.params = options.params.set('provider', String(args.provider));
    }
    return this.sendRequest<models.RskAuthURL>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listAuthorContributions(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskAuthorContributionList> {
    const path = `/api/author/contribution`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('filter' in args) {
      options.params = options.params.set('filter', String(args.filter));
    }
    if ('sortField' in args) {
      options.params = options.params.set('sortField', String(args.sortField));
    }
    if ('sortDirection' in args) {
      options.params = options.params.set('sortDirection', String(args.sortDirection));
    }
    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    if ('pageSize' in args) {
      options.params = options.params.set('pageSize', String(args.pageSize));
    }
    return this.sendRequest<models.RskAuthorContributionList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listAuthorRanks(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskAuthorRankList> {
    const path = `/api/author/ranks`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('filter' in args) {
      options.params = options.params.set('filter', String(args.filter));
    }
    if ('sortField' in args) {
      options.params = options.params.set('sortField', String(args.sortField));
    }
    if ('sortDirection' in args) {
      options.params = options.params.set('sortDirection', String(args.sortDirection));
    }
    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    if ('pageSize' in args) {
      options.params = options.params.set('pageSize', String(args.pageSize));
    }
    return this.sendRequest<models.RskAuthorRankList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listChangelogs(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChangelogList> {
    const path = `/api/changelog`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('filter' in args) {
      options.params = options.params.set('filter', String(args.filter));
    }
    if ('sortField' in args) {
      options.params = options.params.set('sortField', String(args.sortField));
    }
    if ('sortDirection' in args) {
      options.params = options.params.set('sortDirection', String(args.sortDirection));
    }
    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    if ('pageSize' in args) {
      options.params = options.params.set('pageSize', String(args.pageSize));
    }
    return this.sendRequest<models.RskChangelogList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listIncomingDonations(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskIncomingDonationList> {
    const path = `/api/donations`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskIncomingDonationList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getMetadata(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskMetadata> {
    const path = `/api/metadata`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskMetadata>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listPendingRewards(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskPendingRewardList> {
    const path = `/api/rewards`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskPendingRewardList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listClaimedRewards(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskClaimedRewardList> {
    const path = `/api/rewards/claimed`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskClaimedRewardList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getDonationStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskDonationStats> {
    const path = `/api/rewards/stats`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskDonationStats>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  claimReward(
    args: {
      id: string,
      body: models.RskClaimRewardRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    const path = `/api/rewards/${args.id}/claim`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<object>('PATCH', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listDonationRecipients(
    args: {
      rewardId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskDonationRecipientList> {
    const path = `/api/rewards/${args.rewardId}/recipients`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskDonationRecipientList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  search(
    args: {
      query?: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchResultList> {
    const path = `/api/search`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('query' in args) {
      options.params = options.params.set('query', String(args.query));
    }
    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    return this.sendRequest<models.RskSearchResultList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  predictSearchTerm(
    args: {
      prefix?: string,
      maxPredictions?: number,
      query?: string,
      exact?: boolean,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchTermPredictions> {
    const path = `/api/search/predict-terms`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('prefix' in args) {
      options.params = options.params.set('prefix', String(args.prefix));
    }
    if ('maxPredictions' in args) {
      options.params = options.params.set('maxPredictions', String(args.maxPredictions));
    }
    if ('query' in args) {
      options.params = options.params.set('query', String(args.query));
    }
    if ('exact' in args) {
      options.params = options.params.set('exact', String(args.exact));
    }
    return this.sendRequest<models.RskSearchTermPredictions>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getRandomQuote(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskRandomQuote> {
    const path = `/api/search/random-quote`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskRandomQuote>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listSongs(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSongList> {
    const path = `/api/search/songs`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('filter' in args) {
      options.params = options.params.set('filter', String(args.filter));
    }
    if ('sortField' in args) {
      options.params = options.params.set('sortField', String(args.sortField));
    }
    if ('sortDirection' in args) {
      options.params = options.params.set('sortDirection', String(args.sortDirection));
    }
    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    if ('pageSize' in args) {
      options.params = options.params.set('pageSize', String(args.pageSize));
    }
    return this.sendRequest<models.RskSongList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getQuotaSummary(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskQuotas> {
    const path = `/api/status/quotas`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskQuotas>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listTranscripts(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptList> {
    const path = `/api/transcript`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskTranscriptList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listTranscriptChanges(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChangeList> {
    const path = `/api/transcript/change`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('filter' in args) {
      options.params = options.params.set('filter', String(args.filter));
    }
    if ('sortField' in args) {
      options.params = options.params.set('sortField', String(args.sortField));
    }
    if ('sortDirection' in args) {
      options.params = options.params.set('sortDirection', String(args.sortDirection));
    }
    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    if ('pageSize' in args) {
      options.params = options.params.set('pageSize', String(args.pageSize));
    }
    return this.sendRequest<models.RskTranscriptChangeList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getTranscriptChange(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChange> {
    const path = `/api/transcript/change/${args.id}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskTranscriptChange>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  deleteTranscriptChange(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    const path = `/api/transcript/change/${args.id}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<object>('DELETE', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  updateTranscriptChange(
    args: {
      id: string,
      body: models.RskUpdateTranscriptChangeRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChange> {
    const path = `/api/transcript/change/${args.id}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskTranscriptChange>('PATCH', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getTranscriptChangeDiff(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChangeDiff> {
    const path = `/api/transcript/change/${args.id}/diff`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskTranscriptChangeDiff>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  requestTranscriptChangeState(
    args: {
      id: string,
      body: models.RskRequestTranscriptChangeStateRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    const path = `/api/transcript/change/${args.id}/state`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<object>('PATCH', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listChunkedTranscripts(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkedTranscriptList> {
    const path = `/api/transcript/chunked`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskChunkedTranscriptList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getChunkContribution(
    args: {
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution> {
    const path = `/api/transcript/chunked/chunk/contribution/${args.contributionId}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskChunkContribution>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  deleteChunkContribution(
    args: {
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    const path = `/api/transcript/chunked/chunk/contribution/${args.contributionId}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<object>('DELETE', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  updateChunkContribution(
    args: {
      contributionId: string,
      body: models.RskUpdateChunkContributionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution> {
    const path = `/api/transcript/chunked/chunk/contribution/${args.contributionId}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskChunkContribution>('PATCH', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  requestChunkContributionState(
    args: {
      contributionId: string,
      body: models.RskRequestChunkContributionStateRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution> {
    const path = `/api/transcript/chunked/chunk/contribution/${args.contributionId}/state`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskChunkContribution>('PATCH', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listChunkContributions(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContributionList> {
    const path = `/api/transcript/chunked/chunk/contributions`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('filter' in args) {
      options.params = options.params.set('filter', String(args.filter));
    }
    if ('sortField' in args) {
      options.params = options.params.set('sortField', String(args.sortField));
    }
    if ('sortDirection' in args) {
      options.params = options.params.set('sortDirection', String(args.sortDirection));
    }
    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    if ('pageSize' in args) {
      options.params = options.params.set('pageSize', String(args.pageSize));
    }
    return this.sendRequest<models.RskChunkContributionList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  createChunkContribution(
    args: {
      chunkId: string,
      body: models.RskCreateChunkContributionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution> {
    const path = `/api/transcript/chunked/chunk/${args.chunkId}/contribution`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskChunkContribution>('POST', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getTranscriptChunk(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunk> {
    const path = `/api/transcript/chunked/chunk/${args.id}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskChunk>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listTranscriptChunks(
    args: {
      chunkedTranscriptId: string,
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChunkList> {
    const path = `/api/transcript/chunked/${args.chunkedTranscriptId}/chunks`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('filter' in args) {
      options.params = options.params.set('filter', String(args.filter));
    }
    if ('sortField' in args) {
      options.params = options.params.set('sortField', String(args.sortField));
    }
    if ('sortDirection' in args) {
      options.params = options.params.set('sortDirection', String(args.sortDirection));
    }
    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    if ('pageSize' in args) {
      options.params = options.params.set('pageSize', String(args.pageSize));
    }
    return this.sendRequest<models.RskTranscriptChunkList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getTranscript(
    args: {
      epid: string,
      withRaw?: boolean,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscript> {
    const path = `/api/transcript/${args.epid}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('withRaw' in args) {
      options.params = options.params.set('withRaw', String(args.withRaw));
    }
    return this.sendRequest<models.RskTranscript>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  createTranscriptChange(
    args: {
      epid: string,
      body: models.RskCreateTranscriptChangeRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChange> {
    const path = `/api/transcript/${args.epid}/change`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskTranscriptChange>('POST', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getTranscriptDialog(
    args: {
      epid: string,
      pos: number,
      numContextLines?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptDialog> {
    const path = `/api/transcript/${args.epid}/dialog/${args.pos}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('numContextLines' in args) {
      options.params = options.params.set('numContextLines', String(args.numContextLines));
    }
    return this.sendRequest<models.RskTranscriptDialog>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getChunkedTranscriptChunkStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkStats> {
    const path = `/api/transcripts/chunked/chunk-stats`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskChunkStats>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listNotifications(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskNotificationsList> {
    const path = `/api/user/notifications`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('filter' in args) {
      options.params = options.params.set('filter', String(args.filter));
    }
    if ('sortField' in args) {
      options.params = options.params.set('sortField', String(args.sortField));
    }
    if ('sortDirection' in args) {
      options.params = options.params.set('sortDirection', String(args.sortDirection));
    }
    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    if ('pageSize' in args) {
      options.params = options.params.set('pageSize', String(args.pageSize));
    }
    return this.sendRequest<models.RskNotificationsList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  markNotificationsRead(
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    const path = `/api/user/notifications/mark-all`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<object>('POST', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listFieldValues(
    args: {
      field: string,
      prefix?: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskFieldValueList> {
    const path = `/api/values/${args.field}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('prefix' in args) {
      options.params = options.params.set('prefix', String(args.prefix));
    }
    return this.sendRequest<models.RskFieldValueList>('GET', path, options);
  }

  private sendRequest<T>(method: string, path: string, options: HttpOptions, body?: any): Observable<T> {
    switch (method) {
      case 'DELETE':
        return this.http.delete<T>(`${this.domain}${path}`, options);
      case 'GET':
        return this.http.get<T>(`${this.domain}${path}`, options);
      case 'HEAD':
        return this.http.head<T>(`${this.domain}${path}`, options);
      case 'OPTIONS':
        return this.http.options<T>(`${this.domain}${path}`, options);
      case 'PATCH':
        return this.http.patch<T>(`${this.domain}${path}`, body, options);
      case 'POST':
        return this.http.post<T>(`${this.domain}${path}`, body, options);
      case 'PUT':
        return this.http.put<T>(`${this.domain}${path}`, body, options);
      default:
        console.error(`Unsupported request: ${method}`);
        return throwError(`Unsupported request: ${method}`);
    }
  }
}
