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
      body: models.RsksearchSubmitDialogCorrectionRequest,
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
  searchServiceListTscripts(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceListAuthorContributions(
    args: {
      authorId: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContributionList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceCreateChunkContribution(
    args: {
      chunkId: string,
      body: models.RsksearchCreateChunkContributionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContribution>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetChunkContribution(
    args: {
      chunkId: string,
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContribution>;

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
  ): Observable<models.RsksearchChunkContribution>;

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
  ): Observable<models.RsksearchChunkContribution>;

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
  searchServiceListTscriptChunkContributions(
    args: {
      tscriptId: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptChunkContributionList>;

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
