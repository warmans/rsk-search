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
  getRedditAuthURL(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskRedditAuthURL> {
    const path = `/api/auth/reddit-url`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskRedditAuthURL>('GET', path, options);
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
    const path = `/api/contrib/chunk`;
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
    const path = `/api/contrib/chunk/${args.chunkId}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskChunkContribution>('POST', path, options, JSON.stringify(args.body));
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
    const path = `/api/contrib/chunk/${args.contributionId}`;
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
    const path = `/api/contrib/chunk/${args.contributionId}`;
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
    const path = `/api/contrib/chunk/${args.contributionId}`;
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
    const path = `/api/contrib/chunk/${args.contributionId}/state`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskChunkContribution>('PATCH', path, options, JSON.stringify(args.body));
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
    const path = `/api/contrib/transcript/change/${args.id}`;
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
    const path = `/api/contrib/transcript/change/${args.id}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskTranscriptChange>('PATCH', path, options, JSON.stringify(args.body));
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
    const path = `/api/contrib/transcript/change/${args.id}/state`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<object>('PATCH', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listTranscriptChanges(
    args: {
      transcriptId: string,
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChangeList> {
    const path = `/api/contrib/transcript/${args.transcriptId}/change`;
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
  createTranscriptChange(
    args: {
      transcriptId: string,
      body: models.RskCreateTranscriptChangeRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChange> {
    const path = `/api/contrib/transcript/${args.transcriptId}/change`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskTranscriptChange>('POST', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getSearchMetadata(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchMetadata> {
    const path = `/api/metadata`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskSearchMetadata>('GET', path, options);
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
  getTranscript(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscript> {
    const path = `/api/transcript/${args.id}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskTranscript>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listTscripts(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTscriptList> {
    const path = `/api/tscript`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskTscriptList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getAuthorLeaderboard(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskAuthorLeaderboard> {
    const path = `/api/tscript/author/leaderboard`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskAuthorLeaderboard>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getChunk(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunk> {
    const path = `/api/tscript/chunk/${args.id}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskChunk>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getChunkStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkStats> {
    const path = `/api/tscript/stats`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RskChunkStats>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listChunks(
    args: {
      tscriptId: string,
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkList> {
    const path = `/api/tscript/${args.tscriptId}/chunk`;
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
    return this.sendRequest<models.RskChunkList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getTscriptTimeline(
    args: {
      tscriptId: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTscriptTimeline> {
    const path = `/api/tscript/${args.tscriptId}/timeline`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    return this.sendRequest<models.RskTscriptTimeline>('GET', path, options);
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
