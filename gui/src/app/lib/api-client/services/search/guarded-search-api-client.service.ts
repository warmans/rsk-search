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

  searchServiceGetRedditAuthURL(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchRedditAuthURL> {
    return super.searchServiceGetRedditAuthURL(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchRedditAuthURL(res) || console.error(`TypeGuard for response 'RsksearchRedditAuthURL' caught inconsistency.`, res)));
  }

  searchServiceListEpisodes(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchEpisodeList> {
    return super.searchServiceListEpisodes(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchEpisodeList(res) || console.error(`TypeGuard for response 'RsksearchEpisodeList' caught inconsistency.`, res)));
  }

  searchServiceSubmitDialogCorrection(
    args: {
      episodeId: string,
      id: string,
      body: models.RsksearchSubmitDialogCorrectionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    return super.searchServiceSubmitDialogCorrection(args, requestHttpOptions)
      .pipe(tap((res: any) => typeof res === 'object' || console.error(`TypeGuard for response 'object' caught inconsistency.`, res)));
  }

  searchServiceGetEpisode(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchEpisode> {
    return super.searchServiceGetEpisode(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchEpisode(res) || console.error(`TypeGuard for response 'RsksearchEpisode' caught inconsistency.`, res)));
  }

  searchServiceGetSearchMetadata(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchMetadata> {
    return super.searchServiceGetSearchMetadata(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskSearchMetadata(res) || console.error(`TypeGuard for response 'RskSearchMetadata' caught inconsistency.`, res)));
  }

  searchServiceListPendingRewards(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchPendingRewardList> {
    return super.searchServiceListPendingRewards(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchPendingRewardList(res) || console.error(`TypeGuard for response 'RsksearchPendingRewardList' caught inconsistency.`, res)));
  }

  searchServiceListClaimedRewards(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchClaimedRewardList> {
    return super.searchServiceListClaimedRewards(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchClaimedRewardList(res) || console.error(`TypeGuard for response 'RsksearchClaimedRewardList' caught inconsistency.`, res)));
  }

  searchServiceClaimReward(
    args: {
      id: string,
      body: models.RsksearchClaimRewardRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    return super.searchServiceClaimReward(args, requestHttpOptions)
      .pipe(tap((res: any) => typeof res === 'object' || console.error(`TypeGuard for response 'object' caught inconsistency.`, res)));
  }

  searchServiceListDonationRecipients(
    args: {
      rewardId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchDonationRecipientList> {
    return super.searchServiceListDonationRecipients(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchDonationRecipientList(res) || console.error(`TypeGuard for response 'RsksearchDonationRecipientList' caught inconsistency.`, res)));
  }

  searchServiceSearch(
    args: {
      query?: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchResultList> {
    return super.searchServiceSearch(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskSearchResultList(res) || console.error(`TypeGuard for response 'RskSearchResultList' caught inconsistency.`, res)));
  }

  searchServiceListTscripts(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptList> {
    return super.searchServiceListTscripts(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchTscriptList(res) || console.error(`TypeGuard for response 'RsksearchTscriptList' caught inconsistency.`, res)));
  }

  searchServiceGetAuthorLeaderboard(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchAuthorLeaderboard> {
    return super.searchServiceGetAuthorLeaderboard(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchAuthorLeaderboard(res) || console.error(`TypeGuard for response 'RsksearchAuthorLeaderboard' caught inconsistency.`, res)));
  }

  searchServiceListAuthorContributions(
    args: {
      authorId: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContributionList> {
    return super.searchServiceListAuthorContributions(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchChunkContributionList(res) || console.error(`TypeGuard for response 'RsksearchChunkContributionList' caught inconsistency.`, res)));
  }

  searchServiceCreateChunkContribution(
    args: {
      chunkId: string,
      body: models.RsksearchCreateChunkContributionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContribution> {
    return super.searchServiceCreateChunkContribution(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchChunkContribution(res) || console.error(`TypeGuard for response 'RsksearchChunkContribution' caught inconsistency.`, res)));
  }

  searchServiceGetChunkContribution(
    args: {
      chunkId: string,
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContribution> {
    return super.searchServiceGetChunkContribution(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchChunkContribution(res) || console.error(`TypeGuard for response 'RsksearchChunkContribution' caught inconsistency.`, res)));
  }

  searchServiceDiscardDraftContribution(
    args: {
      chunkId: string,
      contributionId: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    return super.searchServiceDiscardDraftContribution(args, requestHttpOptions)
      .pipe(tap((res: any) => typeof res === 'object' || console.error(`TypeGuard for response 'object' caught inconsistency.`, res)));
  }

  searchServiceUpdateChunkContribution(
    args: {
      chunkId: string,
      contributionId: string,
      body: models.RsksearchUpdateChunkContributionRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContribution> {
    return super.searchServiceUpdateChunkContribution(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchChunkContribution(res) || console.error(`TypeGuard for response 'RsksearchChunkContribution' caught inconsistency.`, res)));
  }

  searchServiceRequestChunkContributionState(
    args: {
      chunkId: string,
      contributionId: string,
      body: models.RsksearchRequestChunkContributionStateRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkContribution> {
    return super.searchServiceRequestChunkContributionState(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchChunkContribution(res) || console.error(`TypeGuard for response 'RsksearchChunkContribution' caught inconsistency.`, res)));
  }

  searchServiceGetTscriptChunk(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptChunk> {
    return super.searchServiceGetTscriptChunk(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchTscriptChunk(res) || console.error(`TypeGuard for response 'RsksearchTscriptChunk' caught inconsistency.`, res)));
  }

  searchServiceGetTscriptChunkStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchChunkStats> {
    return super.searchServiceGetTscriptChunkStats(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchChunkStats(res) || console.error(`TypeGuard for response 'RsksearchChunkStats' caught inconsistency.`, res)));
  }

  searchServiceListTscriptChunkContributions(
    args: {
      tscriptId: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptChunkContributionList> {
    return super.searchServiceListTscriptChunkContributions(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchTscriptChunkContributionList(res) || console.error(`TypeGuard for response 'RsksearchTscriptChunkContributionList' caught inconsistency.`, res)));
  }

  searchServiceGetTscriptTimeline(
    args: {
      tscriptId: string,
      page?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchTscriptTimeline> {
    return super.searchServiceGetTscriptTimeline(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchTscriptTimeline(res) || console.error(`TypeGuard for response 'RsksearchTscriptTimeline' caught inconsistency.`, res)));
  }

  searchServiceListFieldValues(
    args: {
      field: string,
      prefix?: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchFieldValueList> {
    return super.searchServiceListFieldValues(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRsksearchFieldValueList(res) || console.error(`TypeGuard for response 'RsksearchFieldValueList' caught inconsistency.`, res)));
  }

}
