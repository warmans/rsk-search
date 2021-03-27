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
  searchServiceGetRedditAuthURL(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchRedditAuthURL> {
    const path = `/api/auth/reddit-url`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RsksearchRedditAuthURL>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceListEpisodes(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchEpisodeList> {
    const path = `/api/episode`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RsksearchEpisodeList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceSubmitDialogCorrection(
    args: {
      episodeId: string,
      id: string,
      body: models.RsksearchSubmitDialogCorrectionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    const path = `/api/episode/${args.episodeId}/dialog/${args.id}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<object>('POST', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetEpisode(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchEpisode> {
    const path = `/api/episode/${args.id}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RsksearchEpisode>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetSearchMetadata(
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
  searchServiceSearch(
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
  searchServiceListTscripts(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptList> {
    const path = `/api/tscript`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RsksearchTscriptList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceListAuthorContributions(
    args: {
      authorId: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContributionList> {
    const path = `/api/tscript/author/${args.authorId}/contrib`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    return this.sendRequest<models.RsksearchChunkContributionList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceCreateChunkContribution(
    args: {
      chunkId: string,
      body: models.RsksearchCreateChunkContributionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContribution> {
    const path = `/api/tscript/chunk/${args.chunkId}/contrib`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RsksearchChunkContribution>('PATCH', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetChunkContribution(
    args: {
      chunkId: string,
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContribution> {
    const path = `/api/tscript/chunk/${args.chunkId}/contrib/${args.contributionId}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RsksearchChunkContribution>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceUpdateChunkContribution(
    args: {
      chunkId: string,
      contributionId: string,
      body: models.RsksearchUpdateChunkContributionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContribution> {
    const path = `/api/tscript/chunk/${args.chunkId}/contrib/${args.contributionId}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RsksearchChunkContribution>('PATCH', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceRequestChunkContributionState(
    args: {
      chunkId: string,
      contributionId: string,
      body: models.RsksearchRequestChunkContributionStateRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContribution> {
    const path = `/api/tscript/chunk/${args.chunkId}/contrib/${args.contributionId}/state`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RsksearchChunkContribution>('PATCH', path, options, JSON.stringify(args.body));
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetTscriptChunk(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptChunk> {
    const path = `/api/tscript/chunk/${args.id}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RsksearchTscriptChunk>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetTscriptChunkStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkStats> {
    const path = `/api/tscript/stats`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    return this.sendRequest<models.RsksearchChunkStats>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceListTscriptChunkContributions(
    args: {
      tscriptId: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptChunkContributionList> {
    const path = `/api/tscript/${args.tscriptId}/contrib`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('page' in args) {
      options.params = options.params.set('page', String(args.page));
    }
    return this.sendRequest<models.RsksearchTscriptChunkContributionList>('GET', path, options);
  }

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceListFieldValues(
    args: {
      field: string,
      prefix?: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchFieldValueList> {
    const path = `/api/values/${args.field}`;
    const options: APIHttpOptions = {
      ...this.options,
      ...requestHttpOptions,
    };

    if ('prefix' in args) {
      options.params = options.params.set('prefix', String(args.prefix));
    }
    return this.sendRequest<models.RsksearchFieldValueList>('GET', path, options);
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
