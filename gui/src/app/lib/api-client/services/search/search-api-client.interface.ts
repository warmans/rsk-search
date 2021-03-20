/* tslint:disable */

import { Observable } from 'rxjs';
import { HttpOptions } from './';
import * as models from '../../models';

export interface SearchAPIClientInterface {

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetRedditAuthURL(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchRedditAuthURL>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceAuthorizeRedditToken(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchToken>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceListEpisodes(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchEpisodeList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceSubmitDialogCorrection(
    args: {
      episodeId: string,
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetEpisode(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchEpisode>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetSearchMetadata(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchMetadata>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceSearch(
    args: {
      query?: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchResultList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceSubmitTscriptChunk(
    args: {
      chunkId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceListTscriptChunkSubmissions(
    args: {
      chunkId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkSubmissionList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetTscriptChunk(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptChunk>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetTscriptChunkStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkStats>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceListFieldValues(
    args: {
      field: string,
      prefix?: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchFieldValueList>;

}
