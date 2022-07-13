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
  listAuthorContributions(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskAuthorContributionList>;

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
  ): Observable<models.RskAuthorRankList>;

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
  ): Observable<models.RskChangelogList>;

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
  ): Observable<models.RskChunkContributionList>;

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
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  deleteChunkContribution(
    args: {
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  updateChunkContribution(
    args: {
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
      contributionId: string,
      body: models.RskRequestChunkContributionStateRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution>;

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
  ): Observable<models.RskTranscriptChangeList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getTranscriptChange(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChange>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  deleteTranscriptChange(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  updateTranscriptChange(
    args: {
      id: string,
      body: models.RskUpdateTranscriptChangeRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChange>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getTranscriptChangeDiff(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChangeDiff>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  requestTranscriptChangeState(
    args: {
      id: string,
      body: models.RskRequestTranscriptChangeStateRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  createTranscriptChange(
    args: {
      epid: string,
      body: models.RskCreateTranscriptChangeRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChange>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getMetadata(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskMetadata>;

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
  listTranscripts(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getTranscript(
    args: {
      epid: string,
      withRaw?: boolean,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscript>;

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
  getChunk(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunk>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  createTscriptImport(
    args: {
      body: models.RskCreateTscriptImportRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTscriptImport>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  getChunkStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkStats>;

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
  ): Observable<models.RskChunkList>;

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
