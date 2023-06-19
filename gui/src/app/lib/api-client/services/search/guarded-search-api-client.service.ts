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

  createTscriptImport(
    args: {
      body: models.RskCreateTscriptImportRequest,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTscriptImport> {
    return super.createTscriptImport(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTscriptImport(res) || console.error(`TypeGuard for response 'RskTscriptImport' caught inconsistency.`, res)));
  }

  listTscriptImports(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTscriptImportList> {
    return super.listTscriptImports(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTscriptImportList(res) || console.error(`TypeGuard for response 'RskTscriptImportList' caught inconsistency.`, res)));
  }

  deleteTscript(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    return super.deleteTscript(args, requestHttpOptions)
      .pipe(tap((res: any) => typeof res === 'object' || console.error(`TypeGuard for response 'object' caught inconsistency.`, res)));
  }

  getRedditAuthURL(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskRedditAuthURL> {
    return super.getRedditAuthURL(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskRedditAuthURL(res) || console.error(`TypeGuard for response 'RskRedditAuthURL' caught inconsistency.`, res)));
  }

  listAuthorContributions(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskAuthorContributionList> {
    return super.listAuthorContributions(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskAuthorContributionList(res) || console.error(`TypeGuard for response 'RskAuthorContributionList' caught inconsistency.`, res)));
  }

  listAuthorRanks(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskAuthorRankList> {
    return super.listAuthorRanks(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskAuthorRankList(res) || console.error(`TypeGuard for response 'RskAuthorRankList' caught inconsistency.`, res)));
  }

  listChangelogs(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChangelogList> {
    return super.listChangelogs(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChangelogList(res) || console.error(`TypeGuard for response 'RskChangelogList' caught inconsistency.`, res)));
  }

  listIncomingDonations(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskIncomingDonationList> {
    return super.listIncomingDonations(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskIncomingDonationList(res) || console.error(`TypeGuard for response 'RskIncomingDonationList' caught inconsistency.`, res)));
  }

  getMetadata(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskMetadata> {
    return super.getMetadata(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskMetadata(res) || console.error(`TypeGuard for response 'RskMetadata' caught inconsistency.`, res)));
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

  getDonationStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskDonationStats> {
    return super.getDonationStats(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskDonationStats(res) || console.error(`TypeGuard for response 'RskDonationStats' caught inconsistency.`, res)));
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

  predictSearchTerm(
    args: {
      prefix?: string,
      maxPredictions?: number,
      query?: string,
      exact?: boolean,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchTermPredictions> {
    return super.predictSearchTerm(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskSearchTermPredictions(res) || console.error(`TypeGuard for response 'RskSearchTermPredictions' caught inconsistency.`, res)));
  }

  getQuotaSummary(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskQuotas> {
    return super.getQuotaSummary(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskQuotas(res) || console.error(`TypeGuard for response 'RskQuotas' caught inconsistency.`, res)));
  }

  listTranscripts(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptList> {
    return super.listTranscripts(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTranscriptList(res) || console.error(`TypeGuard for response 'RskTranscriptList' caught inconsistency.`, res)));
  }

  listTranscriptChanges(
    args: {
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

  getTranscriptChangeDiff(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChangeDiff> {
    return super.getTranscriptChangeDiff(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTranscriptChangeDiff(res) || console.error(`TypeGuard for response 'RskTranscriptChangeDiff' caught inconsistency.`, res)));
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

  listChunkedTranscripts(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkedTranscriptList> {
    return super.listChunkedTranscripts(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunkedTranscriptList(res) || console.error(`TypeGuard for response 'RskChunkedTranscriptList' caught inconsistency.`, res)));
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

  getTranscriptChunk(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunk> {
    return super.getTranscriptChunk(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunk(res) || console.error(`TypeGuard for response 'RskChunk' caught inconsistency.`, res)));
  }

  listTranscriptChunks(
    args: {
      chunkedTranscriptId: string,
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskTranscriptChunkList> {
    return super.listTranscriptChunks(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskTranscriptChunkList(res) || console.error(`TypeGuard for response 'RskTranscriptChunkList' caught inconsistency.`, res)));
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

  getChunkedTranscriptChunkStats(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskChunkStats> {
    return super.getChunkedTranscriptChunkStats(requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskChunkStats(res) || console.error(`TypeGuard for response 'RskChunkStats' caught inconsistency.`, res)));
  }

  listNotifications(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskNotificationsList> {
    return super.listNotifications(args, requestHttpOptions)
      .pipe(tap((res: any) => guards.isRskNotificationsList(res) || console.error(`TypeGuard for response 'RskNotificationsList' caught inconsistency.`, res)));
  }

  markNotificationsRead(
    requestHttpOptions?: HttpOptions
  ): Observable<object> {
    return super.markNotificationsRead(requestHttpOptions)
      .pipe(tap((res: any) => typeof res === 'object' || console.error(`TypeGuard for response 'object' caught inconsistency.`, res)));
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
