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

  listEpisodes(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskEpisodeList> {
    return super.listEpisodes(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskEpisodeList(res) || console.error(`TypeGuard for response 'RskEpisodeList' caught inconsistency.`, res)));
  }

  getEpisode(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskEpisode> {
    return super.getEpisode(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskEpisode(res) || console.error(`TypeGuard for response 'RskEpisode' caught inconsistency.`, res)));
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

  getChunk(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunk> {
    return super.getChunk(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunk(res) || console.error(`TypeGuard for response 'RskChunk' caught inconsistency.`, res)));
  }

  listContributions(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskContributionList> {
    return super.listContributions(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskContributionList(res) || console.error(`TypeGuard for response 'RskContributionList' caught inconsistency.`, res)));
  }

  getContribution(
    args: {
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution> {
    return super.getContribution(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunkContribution(res) || console.error(`TypeGuard for response 'RskChunkContribution' caught inconsistency.`, res)));
  }

  deleteContribution(
    args: {
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    return super.deleteContribution(args, requestHttpOptions)
      .pipe(tap((res: any) => typeof res === 'object' || console.error(`TypeGuard for response 'object' caught inconsistency.`, res)));
  }

  updateContribution(
    args: {
      contributionId: string,
      body: models.RskUpdateContributionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution> {
    return super.updateContribution(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunkContribution(res) || console.error(`TypeGuard for response 'RskChunkContribution' caught inconsistency.`, res)));
  }

  requestContributionState(
    args: {
      contributionId: string,
      body: models.RskRequestContributionStateRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkContribution> {
    return super.requestContributionState(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunkContribution(res) || console.error(`TypeGuard for response 'RskChunkContribution' caught inconsistency.`, res)));
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
