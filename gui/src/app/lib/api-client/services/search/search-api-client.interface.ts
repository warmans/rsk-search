/* tslint:disable */

import { Observable } from 'rxjs';
import { HttpOptions } from './';
import * as models from '../../models';

export interface SearchAPIClientInterface {

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getRedditAuthURL(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskRedditAuthURL>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listEpisodes(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskEpisodeList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getEpisode(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskEpisode>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getSearchMetadata(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchMetadata>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listPendingRewards(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskPendingRewardList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listClaimedRewards(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskClaimedRewardList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  claimReward(
    args: {
      id: string,
      body: models.RskClaimRewardRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listDonationRecipients(
    args: {
      rewardId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskDonationRecipientList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  search(
    args: {
      query?: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchResultList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listTscripts(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTscriptList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getAuthorLeaderboard(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskAuthorLeaderboard>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  createChunkContribution(
    args: {
      chunkId: string,
      body: models.RskCreateChunkContributionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getChunkContribution(
    args: {
      chunkId: string,
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  discardDraftContribution(
    args: {
      chunkId: string,
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  updateChunkContribution(
    args: {
      chunkId: string,
      contributionId: string,
      body: models.RskUpdateChunkContributionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  requestChunkContributionState(
    args: {
      chunkId: string,
      contributionId: string,
      body: models.RskRequestChunkContributionStateRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getChunk(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunk>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listContributions(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskContributionList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getChunkStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkStats>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getTscriptTimeline(
    args: {
      tscriptId: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTscriptTimeline>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listFieldValues(
    args: {
      field: string,
      prefix?: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskFieldValueList>;

}
