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

export function isRskAuthor(arg: any): arg is models.RskAuthor {
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

export function isRskAuthorLeaderboard(arg: any): arg is models.RskAuthorLeaderboard {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // authors?: RskAuthorRanking[]
    ( typeof arg.authors === 'undefined' || (Array.isArray(arg.authors) && arg.authors.every((item: unknown) => isRskAuthorRanking(item))) ) &&

  true
  );
  }

export function isRskAuthorRanking(arg: any): arg is models.RskAuthorRanking {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // acceptedContributions?: number
    ( typeof arg.acceptedContributions === 'undefined' || typeof arg.acceptedContributions === 'number' ) &&
    // approver?: boolean
    ( typeof arg.approver === 'undefined' || typeof arg.approver === 'boolean' ) &&
    // author?: RskAuthor
    ( typeof arg.author === 'undefined' || isRskAuthor(arg.author) ) &&
    // awardValue?: number
    ( typeof arg.awardValue === 'undefined' || typeof arg.awardValue === 'number' ) &&

  true
  );
  }

export function isRskChunk(arg: any): arg is models.RskChunk {
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

export function isRskChunkContribution(arg: any): arg is models.RskChunkContribution {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // author?: RskAuthor
    ( typeof arg.author === 'undefined' || isRskAuthor(arg.author) ) &&
    // chunkId?: string
    ( typeof arg.chunkId === 'undefined' || typeof arg.chunkId === 'string' ) &&
    // createdAt?: string
    ( typeof arg.createdAt === 'undefined' || typeof arg.createdAt === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // state?: RskContributionState
    ( typeof arg.state === 'undefined' || isRskContributionState(arg.state) ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }

export function isRskChunkContributionList(arg: any): arg is models.RskChunkContributionList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // contributions?: RskChunkContribution[]
    ( typeof arg.contributions === 'undefined' || (Array.isArray(arg.contributions) && arg.contributions.every((item: unknown) => isRskChunkContribution(item))) ) &&

  true
  );
  }

export function isRskChunkList(arg: any): arg is models.RskChunkList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // chunks?: RskChunk[]
    ( typeof arg.chunks === 'undefined' || (Array.isArray(arg.chunks) && arg.chunks.every((item: unknown) => isRskChunk(item))) ) &&

  true
  );
  }

export function isRskChunkStates(arg: any): arg is models.RskChunkStates {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // states?: RskContributionState[]
    ( typeof arg.states === 'undefined' || (Array.isArray(arg.states) && arg.states.every((item: unknown) => isRskContributionState(item))) ) &&

  true
  );
  }

export function isRskChunkStats(arg: any): arg is models.RskChunkStats {
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

export function isRskClaimedReward(arg: any): arg is models.RskClaimedReward {
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

export function isRskClaimedRewardList(arg: any): arg is models.RskClaimedRewardList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // rewards?: RskClaimedReward[]
    ( typeof arg.rewards === 'undefined' || (Array.isArray(arg.rewards) && arg.rewards.every((item: unknown) => isRskClaimedReward(item))) ) &&

  true
  );
  }

export function isRskClaimRewardRequest(arg: any): arg is models.RskClaimRewardRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // donationArgs?: RskDonationArgs
    ( typeof arg.donationArgs === 'undefined' || isRskDonationArgs(arg.donationArgs) ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&

  true
  );
  }

export function isRskContributionState(arg: any): arg is models.RskContributionState {
  return false
   || arg === models.RskContributionState.STATE_UNDEFINED
   || arg === models.RskContributionState.STATE_REQUEST_APPROVAL
   || arg === models.RskContributionState.STATE_PENDING
   || arg === models.RskContributionState.STATE_APPROVED
   || arg === models.RskContributionState.STATE_REJECTED
  ;
  }

export function isRskCreateChunkContributionRequest(arg: any): arg is models.RskCreateChunkContributionRequest {
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

export function isRskCreateTranscriptChangeRequest(arg: any): arg is models.RskCreateTranscriptChangeRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // epid?: string
    ( typeof arg.epid === 'undefined' || typeof arg.epid === 'string' ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }

export function isRskDialog(arg: any): arg is models.RskDialog {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // actor?: string
    ( typeof arg.actor === 'undefined' || typeof arg.actor === 'string' ) &&
    // content?: string
    ( typeof arg.content === 'undefined' || typeof arg.content === 'string' ) &&
    // contentTags?: { [key: string]: RskTag }
    ( typeof arg.contentTags === 'undefined' || isRskTag(arg.contentTags) ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // isMatchedRow?: boolean
    ( typeof arg.isMatchedRow === 'undefined' || typeof arg.isMatchedRow === 'boolean' ) &&
    // metadata?: { [key: string]: string }
    ( typeof arg.metadata === 'undefined' || typeof arg.metadata === 'string' ) &&
    // notable?: boolean
    ( typeof arg.notable === 'undefined' || typeof arg.notable === 'boolean' ) &&
    // offsetSec?: string
    ( typeof arg.offsetSec === 'undefined' || typeof arg.offsetSec === 'string' ) &&
    // pos?: string
    ( typeof arg.pos === 'undefined' || typeof arg.pos === 'string' ) &&
    // type?: string
    ( typeof arg.type === 'undefined' || typeof arg.type === 'string' ) &&

  true
  );
  }

export function isRskDialogResult(arg: any): arg is models.RskDialogResult {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // lines?: RskDialog[]
    ( typeof arg.lines === 'undefined' || (Array.isArray(arg.lines) && arg.lines.every((item: unknown) => isRskDialog(item))) ) &&
    // score?: number
    ( typeof arg.score === 'undefined' || typeof arg.score === 'number' ) &&

  true
  );
  }

export function isRskDonationArgs(arg: any): arg is models.RskDonationArgs {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // recipient?: string
    ( typeof arg.recipient === 'undefined' || typeof arg.recipient === 'string' ) &&

  true
  );
  }

export function isRskDonationRecipient(arg: any): arg is models.RskDonationRecipient {
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

export function isRskDonationRecipientList(arg: any): arg is models.RskDonationRecipientList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // organizations?: RskDonationRecipient[]
    ( typeof arg.organizations === 'undefined' || (Array.isArray(arg.organizations) && arg.organizations.every((item: unknown) => isRskDonationRecipient(item))) ) &&

  true
  );
  }

export function isRskEditableTranscript(arg: any): arg is models.RskEditableTranscript {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }

export function isRskFieldMeta(arg: any): arg is models.RskFieldMeta {
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

export function isRskFieldValue(arg: any): arg is models.RskFieldValue {
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

export function isRskFieldValueList(arg: any): arg is models.RskFieldValueList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // values?: RskFieldValue[]
    ( typeof arg.values === 'undefined' || (Array.isArray(arg.values) && arg.values.every((item: unknown) => isRskFieldValue(item))) ) &&

  true
  );
  }

export function isRskPendingRewardList(arg: any): arg is models.RskPendingRewardList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // rewards?: RskReward[]
    ( typeof arg.rewards === 'undefined' || (Array.isArray(arg.rewards) && arg.rewards.every((item: unknown) => isRskReward(item))) ) &&

  true
  );
  }

export function isRskRedditAuthURL(arg: any): arg is models.RskRedditAuthURL {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // url?: string
    ( typeof arg.url === 'undefined' || typeof arg.url === 'string' ) &&

  true
  );
  }

export function isRskRequestChunkContributionStateRequest(arg: any): arg is models.RskRequestChunkContributionStateRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // comment?: string
    ( typeof arg.comment === 'undefined' || typeof arg.comment === 'string' ) &&
    // contributionId?: string
    ( typeof arg.contributionId === 'undefined' || typeof arg.contributionId === 'string' ) &&
    // requestState?: RskContributionState
    ( typeof arg.requestState === 'undefined' || isRskContributionState(arg.requestState) ) &&

  true
  );
  }

export function isRskRequestTranscriptChangeStateRequest(arg: any): arg is models.RskRequestTranscriptChangeStateRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // state?: RskContributionState
    ( typeof arg.state === 'undefined' || isRskContributionState(arg.state) ) &&

  true
  );
  }

export function isRskReward(arg: any): arg is models.RskReward {
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

export function isRskSearchMetadata(arg: any): arg is models.RskSearchMetadata {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // fields?: RskFieldMeta[]
    ( typeof arg.fields === 'undefined' || (Array.isArray(arg.fields) && arg.fields.every((item: unknown) => isRskFieldMeta(item))) ) &&

  true
  );
  }

export function isRskSearchResult(arg: any): arg is models.RskSearchResult {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // dialogs?: RskDialogResult[]
    ( typeof arg.dialogs === 'undefined' || (Array.isArray(arg.dialogs) && arg.dialogs.every((item: unknown) => isRskDialogResult(item))) ) &&
    // episode?: RskShortTranscript
    ( typeof arg.episode === 'undefined' || isRskShortTranscript(arg.episode) ) &&

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

export function isRskShortTranscript(arg: any): arg is models.RskShortTranscript {
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

export function isRskSynopsis(arg: any): arg is models.RskSynopsis {
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

export function isRskTag(arg: any): arg is models.RskTag {
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

export function isRskTranscript(arg: any): arg is models.RskTranscript {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // contributors?: string[]
    ( typeof arg.contributors === 'undefined' || (Array.isArray(arg.contributors) && arg.contributors.every((item: unknown) => typeof item === 'string')) ) &&
    // episode?: number
    ( typeof arg.episode === 'undefined' || typeof arg.episode === 'number' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // incomplete?: boolean
    ( typeof arg.incomplete === 'undefined' || typeof arg.incomplete === 'boolean' ) &&
    // metadata?: { [key: string]: string }
    ( typeof arg.metadata === 'undefined' || typeof arg.metadata === 'string' ) &&
    // publication?: string
    ( typeof arg.publication === 'undefined' || typeof arg.publication === 'string' ) &&
    // releaseDate?: string
    ( typeof arg.releaseDate === 'undefined' || typeof arg.releaseDate === 'string' ) &&
    // series?: number
    ( typeof arg.series === 'undefined' || typeof arg.series === 'number' ) &&
    // synopses?: RskSynopsis[]
    ( typeof arg.synopses === 'undefined' || (Array.isArray(arg.synopses) && arg.synopses.every((item: unknown) => isRskSynopsis(item))) ) &&
    // tags?: RskTag[]
    ( typeof arg.tags === 'undefined' || (Array.isArray(arg.tags) && arg.tags.every((item: unknown) => isRskTag(item))) ) &&
    // transcript?: RskDialog[]
    ( typeof arg.transcript === 'undefined' || (Array.isArray(arg.transcript) && arg.transcript.every((item: unknown) => isRskDialog(item))) ) &&

  true
  );
  }

export function isRskTranscriptChange(arg: any): arg is models.RskTranscriptChange {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // author?: RskAuthor
    ( typeof arg.author === 'undefined' || isRskAuthor(arg.author) ) &&
    // createdAt?: string
    ( typeof arg.createdAt === 'undefined' || typeof arg.createdAt === 'string' ) &&
    // diff?: string
    ( typeof arg.diff === 'undefined' || typeof arg.diff === 'string' ) &&
    // episodeId?: string
    ( typeof arg.episodeId === 'undefined' || typeof arg.episodeId === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // state?: RskContributionState
    ( typeof arg.state === 'undefined' || isRskContributionState(arg.state) ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }

export function isRskTranscriptChangeList(arg: any): arg is models.RskTranscriptChangeList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // changes?: RskTranscriptChange[]
    ( typeof arg.changes === 'undefined' || (Array.isArray(arg.changes) && arg.changes.every((item: unknown) => isRskTranscriptChange(item))) ) &&

  true
  );
  }

export function isRskTranscriptList(arg: any): arg is models.RskTranscriptList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // episodes?: RskShortTranscript[]
    ( typeof arg.episodes === 'undefined' || (Array.isArray(arg.episodes) && arg.episodes.every((item: unknown) => isRskShortTranscript(item))) ) &&

  true
  );
  }

export function isRskTscriptList(arg: any): arg is models.RskTscriptList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // tscripts?: RskTscriptStats[]
    ( typeof arg.tscripts === 'undefined' || (Array.isArray(arg.tscripts) && arg.tscripts.every((item: unknown) => isRskTscriptStats(item))) ) &&

  true
  );
  }

export function isRskTscriptStats(arg: any): arg is models.RskTscriptStats {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // chunkContributions?: { [key: string]: RskChunkStates }
    ( typeof arg.chunkContributions === 'undefined' || isRskChunkStates(arg.chunkContributions) ) &&
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

export function isRskTscriptTimeline(arg: any): arg is models.RskTscriptTimeline {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // events?: RskTscriptTimelineEvent[]
    ( typeof arg.events === 'undefined' || (Array.isArray(arg.events) && arg.events.every((item: unknown) => isRskTscriptTimelineEvent(item))) ) &&

  true
  );
  }

export function isRskTscriptTimelineEvent(arg: any): arg is models.RskTscriptTimelineEvent {
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

export function isRskUpdateChunkContributionRequest(arg: any): arg is models.RskUpdateChunkContributionRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // contributionId?: string
    ( typeof arg.contributionId === 'undefined' || typeof arg.contributionId === 'string' ) &&
    // state?: RskContributionState
    ( typeof arg.state === 'undefined' || isRskContributionState(arg.state) ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }

export function isRskUpdateTranscriptChangeRequest(arg: any): arg is models.RskUpdateTranscriptChangeRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }


