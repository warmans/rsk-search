/* tslint:disable */

import * as models from '../models';

/* pre-prepared guards for build in complex types */

function _isBlob(arg: any): arg is Blob {
  return arg != null && typeof arg.size === 'number' && typeof arg.type === 'string' && typeof arg.slice === 'function';
}

export function isFile(arg: any): arg is File {
return arg != null && typeof arg.lastModified === 'number' && typeof arg.name === 'string' && _isBlob(arg);
}

/* generated type guards */

export function isFieldMetaKind(arg: any): arg is models.FieldMetaKind {
  return false
   || arg === models.FieldMetaKind.UNKNOWN
   || arg === models.FieldMetaKind.IDENTIFIER
   || arg === models.FieldMetaKind.KEYWORD
   || arg === models.FieldMetaKind.KEYWORD_LIST
   || arg === models.FieldMetaKind.TEXT
   || arg === models.FieldMetaKind.INT
   || arg === models.FieldMetaKind.FLOAT
   || arg === models.FieldMetaKind.DATE
  ;
  }

export function isRewardKind(arg: any): arg is models.RewardKind {
  return false
   || arg === models.RewardKind.UNKNOWN
   || arg === models.RewardKind.DONATION
  ;
  }

export function isRsksearchAuthor(arg: any): arg is models.RsksearchAuthor {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&

  true
  );
  }

export function isRsksearchAuthorLeaderboard(arg: any): arg is models.RsksearchAuthorLeaderboard {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // authors?: RsksearchAuthorRanking[]
    ( typeof arg.authors === 'undefined' || (Array.isArray(arg.authors) && arg.authors.every((item: unknown) => isRsksearchAuthorRanking(item))) ) &&

  true
  );
  }

export function isRsksearchAuthorRanking(arg: any): arg is models.RsksearchAuthorRanking {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // acceptedContributions?: number
    ( typeof arg.acceptedContributions === 'undefined' || typeof arg.acceptedContributions === 'number' ) &&
    // approver?: boolean
    ( typeof arg.approver === 'undefined' || typeof arg.approver === 'boolean' ) &&
    // authorName?: string
    ( typeof arg.authorName === 'undefined' || typeof arg.authorName === 'string' ) &&

  true
  );
  }

export function isRsksearchChunkContribution(arg: any): arg is models.RsksearchChunkContribution {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // author?: RsksearchAuthor
    ( typeof arg.author === 'undefined' || isRsksearchAuthor(arg.author) ) &&
    // chunkId?: string
    ( typeof arg.chunkId === 'undefined' || typeof arg.chunkId === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // state?: RsksearchContributionState
    ( typeof arg.state === 'undefined' || isRsksearchContributionState(arg.state) ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }

export function isRsksearchChunkContributionList(arg: any): arg is models.RsksearchChunkContributionList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // contributions?: RsksearchShortChunkContribution[]
    ( typeof arg.contributions === 'undefined' || (Array.isArray(arg.contributions) && arg.contributions.every((item: unknown) => isRsksearchShortChunkContribution(item))) ) &&

  true
  );
  }

export function isRsksearchChunkStates(arg: any): arg is models.RsksearchChunkStates {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // states?: RsksearchContributionState[]
    ( typeof arg.states === 'undefined' || (Array.isArray(arg.states) && arg.states.every((item: unknown) => isRsksearchContributionState(item))) ) &&

  true
  );
  }

export function isRsksearchChunkStats(arg: any): arg is models.RsksearchChunkStats {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // numPending?: number
    ( typeof arg.numPending === 'undefined' || typeof arg.numPending === 'number' ) &&
    // numSubmitted?: number
    ( typeof arg.numSubmitted === 'undefined' || typeof arg.numSubmitted === 'number' ) &&
    // suggestedNextChunkId?: string
    ( typeof arg.suggestedNextChunkId === 'undefined' || typeof arg.suggestedNextChunkId === 'string' ) &&

  true
  );
  }

export function isRsksearchClaimedReward(arg: any): arg is models.RsksearchClaimedReward {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // claimAt?: string
    ( typeof arg.claimAt === 'undefined' || typeof arg.claimAt === 'string' ) &&
    // claimCurrency?: string
    ( typeof arg.claimCurrency === 'undefined' || typeof arg.claimCurrency === 'string' ) &&
    // claimDescription?: string
    ( typeof arg.claimDescription === 'undefined' || typeof arg.claimDescription === 'string' ) &&
    // claimKind?: string
    ( typeof arg.claimKind === 'undefined' || typeof arg.claimKind === 'string' ) &&
    // claimValue?: number
    ( typeof arg.claimValue === 'undefined' || typeof arg.claimValue === 'number' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&

  true
  );
  }

export function isRsksearchClaimedRewardList(arg: any): arg is models.RsksearchClaimedRewardList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // rewards?: RsksearchClaimedReward[]
    ( typeof arg.rewards === 'undefined' || (Array.isArray(arg.rewards) && arg.rewards.every((item: unknown) => isRsksearchClaimedReward(item))) ) &&

  true
  );
  }

export function isRsksearchClaimRewardRequest(arg: any): arg is models.RsksearchClaimRewardRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // donationArgs?: RsksearchDonationArgs
    ( typeof arg.donationArgs === 'undefined' || isRsksearchDonationArgs(arg.donationArgs) ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&

  true
  );
  }

export function isRsksearchContributionState(arg: any): arg is models.RsksearchContributionState {
  return false
   || arg === models.RsksearchContributionState.STATE_UNDEFINED
   || arg === models.RsksearchContributionState.STATE_REQUEST_APPROVAL
   || arg === models.RsksearchContributionState.STATE_PENDING
   || arg === models.RsksearchContributionState.STATE_APPROVED
   || arg === models.RsksearchContributionState.STATE_REJECTED
  ;
  }

export function isRsksearchCreateChunkContributionRequest(arg: any): arg is models.RsksearchCreateChunkContributionRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // chunkId?: string
    ( typeof arg.chunkId === 'undefined' || typeof arg.chunkId === 'string' ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }

export function isRsksearchDialog(arg: any): arg is models.RsksearchDialog {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // actor?: string
    ( typeof arg.actor === 'undefined' || typeof arg.actor === 'string' ) &&
    // content?: string
    ( typeof arg.content === 'undefined' || typeof arg.content === 'string' ) &&
    // contentTags?: { [key: string]: RsksearchTag }
    ( typeof arg.contentTags === 'undefined' || isRsksearchTag(arg.contentTags) ) &&
    // contributor?: string
    ( typeof arg.contributor === 'undefined' || typeof arg.contributor === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // isMatchedRow?: boolean
    ( typeof arg.isMatchedRow === 'undefined' || typeof arg.isMatchedRow === 'boolean' ) &&
    // metadata?: { [key: string]: string }
    ( typeof arg.metadata === 'undefined' || typeof arg.metadata === 'string' ) &&
    // notable?: boolean
    ( typeof arg.notable === 'undefined' || typeof arg.notable === 'boolean' ) &&
    // pos?: string
    ( typeof arg.pos === 'undefined' || typeof arg.pos === 'string' ) &&
    // type?: string
    ( typeof arg.type === 'undefined' || typeof arg.type === 'string' ) &&

  true
  );
  }

export function isRsksearchDialogResult(arg: any): arg is models.RsksearchDialogResult {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // lines?: RsksearchDialog[]
    ( typeof arg.lines === 'undefined' || (Array.isArray(arg.lines) && arg.lines.every((item: unknown) => isRsksearchDialog(item))) ) &&
    // score?: number
    ( typeof arg.score === 'undefined' || typeof arg.score === 'number' ) &&

  true
  );
  }

export function isRsksearchDonationArgs(arg: any): arg is models.RsksearchDonationArgs {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // recipient?: string
    ( typeof arg.recipient === 'undefined' || typeof arg.recipient === 'string' ) &&

  true
  );
  }

export function isRsksearchDonationRecipient(arg: any): arg is models.RsksearchDonationRecipient {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // logoUrl?: string
    ( typeof arg.logoUrl === 'undefined' || typeof arg.logoUrl === 'string' ) &&
    // mission?: string
    ( typeof arg.mission === 'undefined' || typeof arg.mission === 'string' ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // ngoId?: string
    ( typeof arg.ngoId === 'undefined' || typeof arg.ngoId === 'string' ) &&
    // quote?: string
    ( typeof arg.quote === 'undefined' || typeof arg.quote === 'string' ) &&
    // url?: string
    ( typeof arg.url === 'undefined' || typeof arg.url === 'string' ) &&

  true
  );
  }

export function isRsksearchDonationRecipientList(arg: any): arg is models.RsksearchDonationRecipientList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // organizations?: RsksearchDonationRecipient[]
    ( typeof arg.organizations === 'undefined' || (Array.isArray(arg.organizations) && arg.organizations.every((item: unknown) => isRsksearchDonationRecipient(item))) ) &&

  true
  );
  }

export function isRsksearchEpisode(arg: any): arg is models.RsksearchEpisode {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // contributors?: string[]
    ( typeof arg.contributors === 'undefined' || (Array.isArray(arg.contributors) && arg.contributors.every((item: unknown) => typeof item === 'string')) ) &&
    // episode?: number
    ( typeof arg.episode === 'undefined' || typeof arg.episode === 'number' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // metadata?: { [key: string]: string }
    ( typeof arg.metadata === 'undefined' || typeof arg.metadata === 'string' ) &&
    // publication?: string
    ( typeof arg.publication === 'undefined' || typeof arg.publication === 'string' ) &&
    // releaseDate?: string
    ( typeof arg.releaseDate === 'undefined' || typeof arg.releaseDate === 'string' ) &&
    // series?: number
    ( typeof arg.series === 'undefined' || typeof arg.series === 'number' ) &&
    // synopses?: RsksearchSynopsis[]
    ( typeof arg.synopses === 'undefined' || (Array.isArray(arg.synopses) && arg.synopses.every((item: unknown) => isRsksearchSynopsis(item))) ) &&
    // tags?: RsksearchTag[]
    ( typeof arg.tags === 'undefined' || (Array.isArray(arg.tags) && arg.tags.every((item: unknown) => isRsksearchTag(item))) ) &&
    // transcript?: RsksearchDialog[]
    ( typeof arg.transcript === 'undefined' || (Array.isArray(arg.transcript) && arg.transcript.every((item: unknown) => isRsksearchDialog(item))) ) &&

  true
  );
  }

export function isRsksearchEpisodeList(arg: any): arg is models.RsksearchEpisodeList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // episodes?: RsksearchShortEpisode[]
    ( typeof arg.episodes === 'undefined' || (Array.isArray(arg.episodes) && arg.episodes.every((item: unknown) => isRsksearchShortEpisode(item))) ) &&

  true
  );
  }

export function isRsksearchFieldMeta(arg: any): arg is models.RsksearchFieldMeta {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // kind?: FieldMetaKind
    ( typeof arg.kind === 'undefined' || isFieldMetaKind(arg.kind) ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&

  true
  );
  }

export function isRsksearchFieldValue(arg: any): arg is models.RsksearchFieldValue {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // count?: number
    ( typeof arg.count === 'undefined' || typeof arg.count === 'number' ) &&
    // value?: string
    ( typeof arg.value === 'undefined' || typeof arg.value === 'string' ) &&

  true
  );
  }

export function isRsksearchFieldValueList(arg: any): arg is models.RsksearchFieldValueList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // values?: RsksearchFieldValue[]
    ( typeof arg.values === 'undefined' || (Array.isArray(arg.values) && arg.values.every((item: unknown) => isRsksearchFieldValue(item))) ) &&

  true
  );
  }

export function isRskSearchMetadata(arg: any): arg is models.RskSearchMetadata {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // fields?: RsksearchFieldMeta[]
    ( typeof arg.fields === 'undefined' || (Array.isArray(arg.fields) && arg.fields.every((item: unknown) => isRsksearchFieldMeta(item))) ) &&

  true
  );
  }

export function isRsksearchPendingRewardList(arg: any): arg is models.RsksearchPendingRewardList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // rewards?: RsksearchReward[]
    ( typeof arg.rewards === 'undefined' || (Array.isArray(arg.rewards) && arg.rewards.every((item: unknown) => isRsksearchReward(item))) ) &&

  true
  );
  }

export function isRsksearchRedditAuthURL(arg: any): arg is models.RsksearchRedditAuthURL {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // url?: string
    ( typeof arg.url === 'undefined' || typeof arg.url === 'string' ) &&

  true
  );
  }

export function isRsksearchRequestChunkContributionStateRequest(arg: any): arg is models.RsksearchRequestChunkContributionStateRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // chunkId?: string
    ( typeof arg.chunkId === 'undefined' || typeof arg.chunkId === 'string' ) &&
    // comment?: string
    ( typeof arg.comment === 'undefined' || typeof arg.comment === 'string' ) &&
    // contributionId?: string
    ( typeof arg.contributionId === 'undefined' || typeof arg.contributionId === 'string' ) &&
    // requestState?: RsksearchContributionState
    ( typeof arg.requestState === 'undefined' || isRsksearchContributionState(arg.requestState) ) &&

  true
  );
  }

export function isRskSearchResult(arg: any): arg is models.RskSearchResult {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // dialogs?: RsksearchDialogResult[]
    ( typeof arg.dialogs === 'undefined' || (Array.isArray(arg.dialogs) && arg.dialogs.every((item: unknown) => isRsksearchDialogResult(item))) ) &&
    // episode?: RsksearchShortEpisode
    ( typeof arg.episode === 'undefined' || isRsksearchShortEpisode(arg.episode) ) &&

  true
  );
  }

export function isRskSearchResultList(arg: any): arg is models.RskSearchResultList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // resultCount?: number
    ( typeof arg.resultCount === 'undefined' || typeof arg.resultCount === 'number' ) &&
    // results?: RskSearchResult[]
    ( typeof arg.results === 'undefined' || (Array.isArray(arg.results) && arg.results.every((item: unknown) => isRskSearchResult(item))) ) &&

  true
  );
  }

export function isRsksearchReward(arg: any): arg is models.RsksearchReward {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // criteria?: string
    ( typeof arg.criteria === 'undefined' || typeof arg.criteria === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // kind?: RewardKind
    ( typeof arg.kind === 'undefined' || isRewardKind(arg.kind) ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // value?: number
    ( typeof arg.value === 'undefined' || typeof arg.value === 'number' ) &&
    // valueCurrency?: string
    ( typeof arg.valueCurrency === 'undefined' || typeof arg.valueCurrency === 'string' ) &&

  true
  );
  }

export function isRsksearchShortChunkContribution(arg: any): arg is models.RsksearchShortChunkContribution {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // authorId?: string
    ( typeof arg.authorId === 'undefined' || typeof arg.authorId === 'string' ) &&
    // chunkId?: string
    ( typeof arg.chunkId === 'undefined' || typeof arg.chunkId === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // state?: RsksearchContributionState
    ( typeof arg.state === 'undefined' || isRsksearchContributionState(arg.state) ) &&

  true
  );
  }

export function isRsksearchShortEpisode(arg: any): arg is models.RsksearchShortEpisode {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // episode?: number
    ( typeof arg.episode === 'undefined' || typeof arg.episode === 'number' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // publication?: string
    ( typeof arg.publication === 'undefined' || typeof arg.publication === 'string' ) &&
    // series?: number
    ( typeof arg.series === 'undefined' || typeof arg.series === 'number' ) &&
    // transcriptAvailable?: boolean
    ( typeof arg.transcriptAvailable === 'undefined' || typeof arg.transcriptAvailable === 'boolean' ) &&

  true
  );
  }

export function isRsksearchSubmitDialogCorrectionRequest(arg: any): arg is models.RsksearchSubmitDialogCorrectionRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // correction?: string
    ( typeof arg.correction === 'undefined' || typeof arg.correction === 'string' ) &&
    // episodeId?: string
    ( typeof arg.episodeId === 'undefined' || typeof arg.episodeId === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&

  true
  );
  }

export function isRsksearchSynopsis(arg: any): arg is models.RsksearchSynopsis {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // description?: string
    ( typeof arg.description === 'undefined' || typeof arg.description === 'string' ) &&
    // endPos?: string
    ( typeof arg.endPos === 'undefined' || typeof arg.endPos === 'string' ) &&
    // startPos?: string
    ( typeof arg.startPos === 'undefined' || typeof arg.startPos === 'string' ) &&

  true
  );
  }

export function isRsksearchTag(arg: any): arg is models.RsksearchTag {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // kind?: string[]
    ( typeof arg.kind === 'undefined' || (Array.isArray(arg.kind) && arg.kind.every((item: unknown) => typeof item === 'string')) ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&

  true
  );
  }

export function isRsksearchTscriptChunk(arg: any): arg is models.RsksearchTscriptChunk {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // audioClipUri?: string
    ( typeof arg.audioClipUri === 'undefined' || typeof arg.audioClipUri === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // numContributions?: number
    ( typeof arg.numContributions === 'undefined' || typeof arg.numContributions === 'number' ) &&
    // raw?: string
    ( typeof arg.raw === 'undefined' || typeof arg.raw === 'string' ) &&
    // tscriptId?: string
    ( typeof arg.tscriptId === 'undefined' || typeof arg.tscriptId === 'string' ) &&

  true
  );
  }

export function isRsksearchTscriptChunkContributionList(arg: any): arg is models.RsksearchTscriptChunkContributionList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // contributions?: RsksearchChunkContribution[]
    ( typeof arg.contributions === 'undefined' || (Array.isArray(arg.contributions) && arg.contributions.every((item: unknown) => isRsksearchChunkContribution(item))) ) &&

  true
  );
  }

export function isRsksearchTscriptList(arg: any): arg is models.RsksearchTscriptList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // tscripts?: RsksearchTscriptStats[]
    ( typeof arg.tscripts === 'undefined' || (Array.isArray(arg.tscripts) && arg.tscripts.every((item: unknown) => isRsksearchTscriptStats(item))) ) &&

  true
  );
  }

export function isRsksearchTscriptStats(arg: any): arg is models.RsksearchTscriptStats {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // chunkContributions?: { [key: string]: RsksearchChunkStates }
    ( typeof arg.chunkContributions === 'undefined' || isRsksearchChunkStates(arg.chunkContributions) ) &&
    // episode?: number
    ( typeof arg.episode === 'undefined' || typeof arg.episode === 'number' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // numApprovedContributions?: number
    ( typeof arg.numApprovedContributions === 'undefined' || typeof arg.numApprovedContributions === 'number' ) &&
    // numChunks?: number
    ( typeof arg.numChunks === 'undefined' || typeof arg.numChunks === 'number' ) &&
    // numContributions?: number
    ( typeof arg.numContributions === 'undefined' || typeof arg.numContributions === 'number' ) &&
    // numPendingContributions?: number
    ( typeof arg.numPendingContributions === 'undefined' || typeof arg.numPendingContributions === 'number' ) &&
    // numRejectedContributions?: number
    ( typeof arg.numRejectedContributions === 'undefined' || typeof arg.numRejectedContributions === 'number' ) &&
    // numRequestApprovalContributions?: number
    ( typeof arg.numRequestApprovalContributions === 'undefined' || typeof arg.numRequestApprovalContributions === 'number' ) &&
    // publication?: string
    ( typeof arg.publication === 'undefined' || typeof arg.publication === 'string' ) &&
    // series?: number
    ( typeof arg.series === 'undefined' || typeof arg.series === 'number' ) &&

  true
  );
  }

export function isRsksearchTscriptTimeline(arg: any): arg is models.RsksearchTscriptTimeline {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // events?: RsksearchTscriptTimelineEvent[]
    ( typeof arg.events === 'undefined' || (Array.isArray(arg.events) && arg.events.every((item: unknown) => isRsksearchTscriptTimelineEvent(item))) ) &&

  true
  );
  }

export function isRsksearchTscriptTimelineEvent(arg: any): arg is models.RsksearchTscriptTimelineEvent {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // what?: string
    ( typeof arg.what === 'undefined' || typeof arg.what === 'string' ) &&
    // when?: string
    ( typeof arg.when === 'undefined' || typeof arg.when === 'string' ) &&
    // who?: string
    ( typeof arg.who === 'undefined' || typeof arg.who === 'string' ) &&

  true
  );
  }

export function isRsksearchUpdateChunkContributionRequest(arg: any): arg is models.RsksearchUpdateChunkContributionRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // chunkId?: string
    ( typeof arg.chunkId === 'undefined' || typeof arg.chunkId === 'string' ) &&
    // contributionId?: string
    ( typeof arg.contributionId === 'undefined' || typeof arg.contributionId === 'string' ) &&
    // state?: RsksearchContributionState
    ( typeof arg.state === 'undefined' || isRsksearchContributionState(arg.state) ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }


