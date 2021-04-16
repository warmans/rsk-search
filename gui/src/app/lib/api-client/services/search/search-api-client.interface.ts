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
  searchServiceListPendingRewards(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchPendingRewardList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceListClaimedRewards(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchClaimedRewardList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceClaimReward(
    args: {
      id: string,
      body: models.RsksearchClaimRewardRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceListDonationRecipients(
    args: {
      rewardId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchDonationRecipientList>;

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
  searchServiceGetAuthorLeaderboard(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchAuthorLeaderboard>;

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
  searchServiceDiscardDraftContribution(
    args: {
      chunkId: string,
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object>;

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
  searchServiceListTscriptContributions(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptContributionList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetTscriptChunkStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkStats>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetTscriptTimeline(
    args: {
      tscriptId: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptTimeline>;

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
