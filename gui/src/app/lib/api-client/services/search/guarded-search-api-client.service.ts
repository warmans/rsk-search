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

  getRedditAuthURL(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskRedditAuthURL> {
    return super.getRedditAuthURL(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskRedditAuthURL(res) || console.error(`TypeGuard for response 'RskRedditAuthURL' caught inconsistency.`, res)));
  }

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
    return super.listChunkContributions(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunkContributionList(res) || console.error(`TypeGuard for response 'RskChunkContributionList' caught inconsistency.`, res)));
  }

  createChunkContribution(
    args: {
      chunkId: string,
      body: models.RskCreateChunkContributionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution> {
    return super.createChunkContribution(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunkContribution(res) || console.error(`TypeGuard for response 'RskChunkContribution' caught inconsistency.`, res)));
  }

  getChunkContribution(
    args: {
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution> {
    return super.getChunkContribution(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunkContribution(res) || console.error(`TypeGuard for response 'RskChunkContribution' caught inconsistency.`, res)));
  }

  deleteChunkContribution(
    args: {
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    return super.deleteChunkContribution(args, requestHttpOptions)
      .pipe(tap((res: any) => typeof res === 'object' || console.error(`TypeGuard for response 'object' caught inconsistency.`, res)));
  }

  updateChunkContribution(
    args: {
      contributionId: string,
      body: models.RskUpdateChunkContributionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution> {
    return super.updateChunkContribution(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunkContribution(res) || console.error(`TypeGuard for response 'RskChunkContribution' caught inconsistency.`, res)));
  }

  requestChunkContributionState(
    args: {
      contributionId: string,
      body: models.RskRequestChunkContributionStateRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution> {
    return super.requestChunkContributionState(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunkContribution(res) || console.error(`TypeGuard for response 'RskChunkContribution' caught inconsistency.`, res)));
  }

  getTranscriptChange(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChange> {
    return super.getTranscriptChange(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTranscriptChange(res) || console.error(`TypeGuard for response 'RskTranscriptChange' caught inconsistency.`, res)));
  }

  deleteTranscriptChange(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    return super.deleteTranscriptChange(args, requestHttpOptions)
      .pipe(tap((res: any) => typeof res === 'object' || console.error(`TypeGuard for response 'object' caught inconsistency.`, res)));
  }

  updateTranscriptChange(
    args: {
      id: string,
      body: models.RskUpdateTranscriptChangeRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChange> {
    return super.updateTranscriptChange(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTranscriptChange(res) || console.error(`TypeGuard for response 'RskTranscriptChange' caught inconsistency.`, res)));
  }

  requestTranscriptChangeState(
    args: {
      id: string,
      body: models.RskRequestTranscriptChangeStateRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    return super.requestTranscriptChangeState(args, requestHttpOptions)
      .pipe(tap((res: any) => typeof res === 'object' || console.error(`TypeGuard for response 'object' caught inconsistency.`, res)));
  }

  listTranscriptChanges(
    args: {
      epid: string,
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChangeList> {
    return super.listTranscriptChanges(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTranscriptChangeList(res) || console.error(`TypeGuard for response 'RskTranscriptChangeList' caught inconsistency.`, res)));
  }

  createTranscriptChange(
    args: {
      epid: string,
      body: models.RskCreateTranscriptChangeRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChange> {
    return super.createTranscriptChange(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTranscriptChange(res) || console.error(`TypeGuard for response 'RskTranscriptChange' caught inconsistency.`, res)));
  }

  getSearchMetadata(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchMetadata> {
    return super.getSearchMetadata(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskSearchMetadata(res) || console.error(`TypeGuard for response 'RskSearchMetadata' caught inconsistency.`, res)));
  }

  listPendingRewards(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskPendingRewardList> {
    return super.listPendingRewards(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskPendingRewardList(res) || console.error(`TypeGuard for response 'RskPendingRewardList' caught inconsistency.`, res)));
  }

  listClaimedRewards(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskClaimedRewardList> {
    return super.listClaimedRewards(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskClaimedRewardList(res) || console.error(`TypeGuard for response 'RskClaimedRewardList' caught inconsistency.`, res)));
  }

  claimReward(
    args: {
      id: string,
      body: models.RskClaimRewardRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    return super.claimReward(args, requestHttpOptions)
      .pipe(tap((res: any) => typeof res === 'object' || console.error(`TypeGuard for response 'object' caught inconsistency.`, res)));
  }

  listDonationRecipients(
    args: {
      rewardId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskDonationRecipientList> {
    return super.listDonationRecipients(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskDonationRecipientList(res) || console.error(`TypeGuard for response 'RskDonationRecipientList' caught inconsistency.`, res)));
  }

  search(
    args: {
      query?: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchResultList> {
    return super.search(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskSearchResultList(res) || console.error(`TypeGuard for response 'RskSearchResultList' caught inconsistency.`, res)));
  }

  listTranscripts(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptList> {
    return super.listTranscripts(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTranscriptList(res) || console.error(`TypeGuard for response 'RskTranscriptList' caught inconsistency.`, res)));
  }

  getTranscript(
    args: {
      epid: string,
      withRaw?: boolean,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscript> {
    return super.getTranscript(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTranscript(res) || console.error(`TypeGuard for response 'RskTranscript' caught inconsistency.`, res)));
  }

  listTscripts(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTscriptList> {
    return super.listTscripts(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTscriptList(res) || console.error(`TypeGuard for response 'RskTscriptList' caught inconsistency.`, res)));
  }

  getAuthorLeaderboard(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskAuthorLeaderboard> {
    return super.getAuthorLeaderboard(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskAuthorLeaderboard(res) || console.error(`TypeGuard for response 'RskAuthorLeaderboard' caught inconsistency.`, res)));
  }

  getChunk(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunk> {
    return super.getChunk(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunk(res) || console.error(`TypeGuard for response 'RskChunk' caught inconsistency.`, res)));
  }

  getChunkStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkStats> {
    return super.getChunkStats(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunkStats(res) || console.error(`TypeGuard for response 'RskChunkStats' caught inconsistency.`, res)));
  }

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
    return super.listChunks(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunkList(res) || console.error(`TypeGuard for response 'RskChunkList' caught inconsistency.`, res)));
  }

  getTscriptTimeline(
    args: {
      tscriptId: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTscriptTimeline> {
    return super.getTscriptTimeline(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTscriptTimeline(res) || console.error(`TypeGuard for response 'RskTscriptTimeline' caught inconsistency.`, res)));
  }

  listFieldValues(
    args: {
      field: string,
      prefix?: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskFieldValueList> {
    return super.listFieldValues(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskFieldValueList(res) || console.error(`TypeGuard for response 'RskFieldValueList' caught inconsistency.`, res)));
  }

}
