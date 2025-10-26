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

export function isAuthorContributionType(arg: any): arg is models.AuthorContributionType {
  return false
   || arg === models.AuthorContributionType.CONTRIBUTION_TYPE_UNKNOWN
   || arg === models.AuthorContributionType.CHUNK
   || arg === models.AuthorContributionType.CHANGE
  ;
  }

export function isContributionsServiceClaimRewardBody(arg: any): arg is models.ContributionsServiceClaimRewardBody {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // donationArgs?: RskDonationArgs
    ( typeof arg.donationArgs === 'undefined' || isRskDonationArgs(arg.donationArgs) ) &&

  true
  );
  }

export function isDialogType(arg: any): arg is models.DialogType {
  return false
   || arg === models.DialogType.UNKNOWN
   || arg === models.DialogType.SONG
   || arg === models.DialogType.CHAT
   || arg === models.DialogType.NONE
   || arg === models.DialogType.GAP
  ;
  }

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

export function isNotificationKind(arg: any): arg is models.NotificationKind {
  return false
   || arg === models.NotificationKind.UNDEFINED_KIND
   || arg === models.NotificationKind.CONFIRMATION
   || arg === models.NotificationKind.INFO
   || arg === models.NotificationKind.WARNING
   || arg === models.NotificationKind.SPAM
  ;
  }

export function isRewardKind(arg: any): arg is models.RewardKind {
  return false
   || arg === models.RewardKind.UNKNOWN
   || arg === models.RewardKind.DONATION
  ;
  }

export function isRskArchive(arg: any): arg is models.RskArchive {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // description?: string
    ( typeof arg.description === 'undefined' || typeof arg.description === 'string' ) &&
    // files?: string[]
    ( typeof arg.files === 'undefined' || (Array.isArray(arg.files) && arg.files.every((item: unknown) => typeof item === 'string')) ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // media?: RskFile[]
    ( typeof arg.media === 'undefined' || (Array.isArray(arg.media) && arg.media.every((item: unknown) => isRskFile(item))) ) &&
    // relatedEpisode?: string
    ( typeof arg.relatedEpisode === 'undefined' || typeof arg.relatedEpisode === 'string' ) &&

  true
  );
  }

export function isRskArchiveList(arg: any): arg is models.RskArchiveList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // items?: RskArchive[]
    ( typeof arg.items === 'undefined' || (Array.isArray(arg.items) && arg.items.every((item: unknown) => isRskArchive(item))) ) &&
    // resultCount?: number
    ( typeof arg.resultCount === 'undefined' || typeof arg.resultCount === 'number' ) &&

  true
  );
  }

export function isRskAudioQuality(arg: any): arg is models.RskAudioQuality {
  return false
   || arg === models.RskAudioQuality.AUDIO_QUALITY_UNKNOWN
   || arg === models.RskAudioQuality.POOR
   || arg === models.RskAudioQuality.AVERAGE
   || arg === models.RskAudioQuality.GOOD
  ;
  }

export function isRskAuthor(arg: any): arg is models.RskAuthor {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // identityIconImg?: string
    ( typeof arg.identityIconImg === 'undefined' || typeof arg.identityIconImg === 'string' ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // oauthProvider?: string
    ( typeof arg.oauthProvider === 'undefined' || typeof arg.oauthProvider === 'string' ) &&
    // supporter?: boolean
    ( typeof arg.supporter === 'undefined' || typeof arg.supporter === 'boolean' ) &&

  true
  );
  }

export function isRskAuthorContribution(arg: any): arg is models.RskAuthorContribution {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // author?: RskAuthor
    ( typeof arg.author === 'undefined' || isRskAuthor(arg.author) ) &&
    // contributionType?: AuthorContributionType
    ( typeof arg.contributionType === 'undefined' || isAuthorContributionType(arg.contributionType) ) &&
    // createdAt?: string
    ( typeof arg.createdAt === 'undefined' || typeof arg.createdAt === 'string' ) &&
    // episodeId?: string
    ( typeof arg.episodeId === 'undefined' || typeof arg.episodeId === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // points?: number
    ( typeof arg.points === 'undefined' || typeof arg.points === 'number' ) &&

  true
  );
  }

export function isRskAuthorContributionList(arg: any): arg is models.RskAuthorContributionList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // contributions?: RskAuthorContribution[]
    ( typeof arg.contributions === 'undefined' || (Array.isArray(arg.contributions) && arg.contributions.every((item: unknown) => isRskAuthorContribution(item))) ) &&

  true
  );
  }

export function isRskAuthorRank(arg: any): arg is models.RskAuthorRank {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // approvedChanges?: number
    ( typeof arg.approvedChanges === 'undefined' || typeof arg.approvedChanges === 'number' ) &&
    // approvedChunks?: number
    ( typeof arg.approvedChunks === 'undefined' || typeof arg.approvedChunks === 'number' ) &&
    // author?: RskAuthor
    ( typeof arg.author === 'undefined' || isRskAuthor(arg.author) ) &&
    // currentRank?: RskRank
    ( typeof arg.currentRank === 'undefined' || isRskRank(arg.currentRank) ) &&
    // nextRank?: RskRank
    ( typeof arg.nextRank === 'undefined' || isRskRank(arg.nextRank) ) &&
    // points?: number
    ( typeof arg.points === 'undefined' || typeof arg.points === 'number' ) &&
    // rewardValueUsd?: number
    ( typeof arg.rewardValueUsd === 'undefined' || typeof arg.rewardValueUsd === 'number' ) &&

  true
  );
  }

export function isRskAuthorRankList(arg: any): arg is models.RskAuthorRankList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // rankings?: RskAuthorRank[]
    ( typeof arg.rankings === 'undefined' || (Array.isArray(arg.rankings) && arg.rankings.every((item: unknown) => isRskAuthorRank(item))) ) &&

  true
  );
  }

export function isRskAuthURL(arg: any): arg is models.RskAuthURL {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // url?: string
    ( typeof arg.url === 'undefined' || typeof arg.url === 'string' ) &&

  true
  );
  }

export function isRskChangelog(arg: any): arg is models.RskChangelog {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // content?: string
    ( typeof arg.content === 'undefined' || typeof arg.content === 'string' ) &&
    // date?: string
    ( typeof arg.date === 'undefined' || typeof arg.date === 'string' ) &&

  true
  );
  }

export function isRskChangelogList(arg: any): arg is models.RskChangelogList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // changelogs?: RskChangelog[]
    ( typeof arg.changelogs === 'undefined' || (Array.isArray(arg.changelogs) && arg.changelogs.every((item: unknown) => isRskChangelog(item))) ) &&

  true
  );
  }

export function isRskChunk(arg: any): arg is models.RskChunk {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // audioClipUri?: string
    ( typeof arg.audioClipUri === 'undefined' || typeof arg.audioClipUri === 'string' ) &&
    // chunkedTranscriptId?: string
    ( typeof arg.chunkedTranscriptId === 'undefined' || typeof arg.chunkedTranscriptId === 'string' ) &&
    // endTimeMs?: number
    ( typeof arg.endTimeMs === 'undefined' || typeof arg.endTimeMs === 'number' ) &&
    // episodeId?: string
    ( typeof arg.episodeId === 'undefined' || typeof arg.episodeId === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // numContributions?: number
    ( typeof arg.numContributions === 'undefined' || typeof arg.numContributions === 'number' ) &&
    // raw?: string
    ( typeof arg.raw === 'undefined' || typeof arg.raw === 'string' ) &&
    // startTimeMs?: number
    ( typeof arg.startTimeMs === 'undefined' || typeof arg.startTimeMs === 'number' ) &&

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
    // stateComment?: string
    ( typeof arg.stateComment === 'undefined' || typeof arg.stateComment === 'string' ) &&
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

export function isRskChunkedTranscriptList(arg: any): arg is models.RskChunkedTranscriptList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // chunked?: RskChunkedTranscriptStats[]
    ( typeof arg.chunked === 'undefined' || (Array.isArray(arg.chunked) && arg.chunked.every((item: unknown) => isRskChunkedTranscriptStats(item))) ) &&

  true
  );
  }

export function isRskChunkedTranscriptStats(arg: any): arg is models.RskChunkedTranscriptStats {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // chunkContributions?: { [key: string]: RskChunkStates }
    ( typeof arg.chunkContributions === 'undefined' || isRskChunkStates(arg.chunkContributions) ) &&
    // episode?: number
    ( typeof arg.episode === 'undefined' || typeof arg.episode === 'number' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
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

export function isRskCommunityProject(arg: any): arg is models.RskCommunityProject {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // content?: string
    ( typeof arg.content === 'undefined' || typeof arg.content === 'string' ) &&
    // createdAt?: string
    ( typeof arg.createdAt === 'undefined' || typeof arg.createdAt === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // summary?: string
    ( typeof arg.summary === 'undefined' || typeof arg.summary === 'string' ) &&
    // url?: string
    ( typeof arg.url === 'undefined' || typeof arg.url === 'string' ) &&

  true
  );
  }

export function isRskCommunityProjectList(arg: any): arg is models.RskCommunityProjectList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // projects?: RskCommunityProject[]
    ( typeof arg.projects === 'undefined' || (Array.isArray(arg.projects) && arg.projects.every((item: unknown) => isRskCommunityProject(item))) ) &&
    // resultCount?: number
    ( typeof arg.resultCount === 'undefined' || typeof arg.resultCount === 'number' ) &&

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

export function isRskCreateTscriptImportRequest(arg: any): arg is models.RskCreateTscriptImportRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // epid?: string
    ( typeof arg.epid === 'undefined' || typeof arg.epid === 'string' ) &&
    // epname?: string
    ( typeof arg.epname === 'undefined' || typeof arg.epname === 'string' ) &&
    // mp3Uri?: string
    ( typeof arg.mp3Uri === 'undefined' || typeof arg.mp3Uri === 'string' ) &&

  true
  );
  }

export function isRskCurrentRadioEpisode(arg: any): arg is models.RskCurrentRadioEpisode {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // shortId?: string
    ( typeof arg.shortId === 'undefined' || typeof arg.shortId === 'string' ) &&
    // startedAt?: string
    ( typeof arg.startedAt === 'undefined' || typeof arg.startedAt === 'string' ) &&

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
    // durationMs?: number
    ( typeof arg.durationMs === 'undefined' || typeof arg.durationMs === 'number' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // isMatchedRow?: boolean
    ( typeof arg.isMatchedRow === 'undefined' || typeof arg.isMatchedRow === 'boolean' ) &&
    // metadata?: { [key: string]: string }
    ( typeof arg.metadata === 'undefined' || typeof arg.metadata === 'string' ) &&
    // notable?: boolean
    ( typeof arg.notable === 'undefined' || typeof arg.notable === 'boolean' ) &&
    // offsetDistance?: number
    ( typeof arg.offsetDistance === 'undefined' || typeof arg.offsetDistance === 'number' ) &&
    // offsetInferred?: boolean
    ( typeof arg.offsetInferred === 'undefined' || typeof arg.offsetInferred === 'boolean' ) &&
    // offsetMs?: number
    ( typeof arg.offsetMs === 'undefined' || typeof arg.offsetMs === 'number' ) &&
    // offsetSec?: string
    ( typeof arg.offsetSec === 'undefined' || typeof arg.offsetSec === 'string' ) &&
    // placeholder?: boolean
    ( typeof arg.placeholder === 'undefined' || typeof arg.placeholder === 'boolean' ) &&
    // pos?: number
    ( typeof arg.pos === 'undefined' || typeof arg.pos === 'number' ) &&
    // type?: DialogType
    ( typeof arg.type === 'undefined' || isDialogType(arg.type) ) &&

  true
  );
  }

export function isRskDialogResult(arg: any): arg is models.RskDialogResult {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // score?: number
    ( typeof arg.score === 'undefined' || typeof arg.score === 'number' ) &&
    // transcript?: RskDialog[]
    ( typeof arg.transcript === 'undefined' || (Array.isArray(arg.transcript) && arg.transcript.every((item: unknown) => isRskDialog(item))) ) &&

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

export function isRskDonationStats(arg: any): arg is models.RskDonationStats {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // stats?: RskRecipientStats[]
    ( typeof arg.stats === 'undefined' || (Array.isArray(arg.stats) && arg.stats.every((item: unknown) => isRskRecipientStats(item))) ) &&

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

export function isRskFile(arg: any): arg is models.RskFile {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // thumbnailName?: string
    ( typeof arg.thumbnailName === 'undefined' || typeof arg.thumbnailName === 'string' ) &&

  true
  );
  }

export function isRskIncomingDonation(arg: any): arg is models.RskIncomingDonation {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // amount?: number
    ( typeof arg.amount === 'undefined' || typeof arg.amount === 'number' ) &&
    // amountCurrency?: string
    ( typeof arg.amountCurrency === 'undefined' || typeof arg.amountCurrency === 'string' ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // note?: string
    ( typeof arg.note === 'undefined' || typeof arg.note === 'string' ) &&

  true
  );
  }

export function isRskIncomingDonationList(arg: any): arg is models.RskIncomingDonationList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // donations?: RskIncomingDonation[]
    ( typeof arg.donations === 'undefined' || (Array.isArray(arg.donations) && arg.donations.every((item: unknown) => isRskIncomingDonation(item))) ) &&

  true
  );
  }

export function isRskMedia(arg: any): arg is models.RskMedia {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // audio?: boolean
    ( typeof arg.audio === 'undefined' || typeof arg.audio === 'boolean' ) &&
    // audioRestricted?: boolean
    ( typeof arg.audioRestricted === 'undefined' || typeof arg.audioRestricted === 'boolean' ) &&
    // video?: boolean
    ( typeof arg.video === 'undefined' || typeof arg.video === 'boolean' ) &&

  true
  );
  }

export function isRskMetadata(arg: any): arg is models.RskMetadata {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // episodeShortIds?: string[]
    ( typeof arg.episodeShortIds === 'undefined' || (Array.isArray(arg.episodeShortIds) && arg.episodeShortIds.every((item: unknown) => typeof item === 'string')) ) &&
    // searchFields?: RskFieldMeta[]
    ( typeof arg.searchFields === 'undefined' || (Array.isArray(arg.searchFields) && arg.searchFields.every((item: unknown) => isRskFieldMeta(item))) ) &&

  true
  );
  }

export function isRskNextRadioEpisode(arg: any): arg is models.RskNextRadioEpisode {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // shortId?: string
    ( typeof arg.shortId === 'undefined' || typeof arg.shortId === 'string' ) &&

  true
  );
  }

export function isRskNotification(arg: any): arg is models.RskNotification {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // clickThoughUrl?: string
    ( typeof arg.clickThoughUrl === 'undefined' || typeof arg.clickThoughUrl === 'string' ) &&
    // createdAt?: string
    ( typeof arg.createdAt === 'undefined' || typeof arg.createdAt === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // kind?: NotificationKind
    ( typeof arg.kind === 'undefined' || isNotificationKind(arg.kind) ) &&
    // message?: string
    ( typeof arg.message === 'undefined' || typeof arg.message === 'string' ) &&
    // readAt?: string
    ( typeof arg.readAt === 'undefined' || typeof arg.readAt === 'string' ) &&

  true
  );
  }

export function isRskNotificationsList(arg: any): arg is models.RskNotificationsList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // notifications?: RskNotification[]
    ( typeof arg.notifications === 'undefined' || (Array.isArray(arg.notifications) && arg.notifications.every((item: unknown) => isRskNotification(item))) ) &&

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

export function isRskPrediction(arg: any): arg is models.RskPrediction {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // actor?: string
    ( typeof arg.actor === 'undefined' || typeof arg.actor === 'string' ) &&
    // epid?: string
    ( typeof arg.epid === 'undefined' || typeof arg.epid === 'string' ) &&
    // fragment?: string
    ( typeof arg.fragment === 'undefined' || typeof arg.fragment === 'string' ) &&
    // line?: string
    ( typeof arg.line === 'undefined' || typeof arg.line === 'string' ) &&
    // pos?: number
    ( typeof arg.pos === 'undefined' || typeof arg.pos === 'number' ) &&

  true
  );
  }

export function isRskPublicationType(arg: any): arg is models.RskPublicationType {
  return false
   || arg === models.RskPublicationType.PUBLICATION_TYPE_UNKNOWN
   || arg === models.RskPublicationType.PUBLICATION_TYPE_RADIO
   || arg === models.RskPublicationType.PUBLICATION_TYPE_PODCAST
   || arg === models.RskPublicationType.PUBLICATION_TYPE_PROMO
   || arg === models.RskPublicationType.PUBLICATION_TYPE_TV
   || arg === models.RskPublicationType.PUBLICATION_TYPE_OTHER
  ;
  }

export function isRskPutRadioStateRequest(arg: any): arg is models.RskPutRadioStateRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // currentEpisode?: RskCurrentRadioEpisode
    ( typeof arg.currentEpisode === 'undefined' || isRskCurrentRadioEpisode(arg.currentEpisode) ) &&
    // currentTimestampMs?: number
    ( typeof arg.currentTimestampMs === 'undefined' || typeof arg.currentTimestampMs === 'number' ) &&

  true
  );
  }

export function isRskQuotas(arg: any): arg is models.RskQuotas {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // bandwidthRemainingMib?: number
    ( typeof arg.bandwidthRemainingMib === 'undefined' || typeof arg.bandwidthRemainingMib === 'number' ) &&
    // bandwidthTotalMib?: number
    ( typeof arg.bandwidthTotalMib === 'undefined' || typeof arg.bandwidthTotalMib === 'number' ) &&

  true
  );
  }

export function isRskRadioState(arg: any): arg is models.RskRadioState {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // currentEpisode?: RskCurrentRadioEpisode
    ( typeof arg.currentEpisode === 'undefined' || isRskCurrentRadioEpisode(arg.currentEpisode) ) &&
    // currentTimestampMs?: number
    ( typeof arg.currentTimestampMs === 'undefined' || typeof arg.currentTimestampMs === 'number' ) &&

  true
  );
  }

export function isRskRandomQuote(arg: any): arg is models.RskRandomQuote {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // actor?: string
    ( typeof arg.actor === 'undefined' || typeof arg.actor === 'string' ) &&
    // epid?: string
    ( typeof arg.epid === 'undefined' || typeof arg.epid === 'string' ) &&
    // pos?: number
    ( typeof arg.pos === 'undefined' || typeof arg.pos === 'number' ) &&
    // quote?: string
    ( typeof arg.quote === 'undefined' || typeof arg.quote === 'string' ) &&

  true
  );
  }

export function isRskRank(arg: any): arg is models.RskRank {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // points?: number
    ( typeof arg.points === 'undefined' || typeof arg.points === 'number' ) &&

  true
  );
  }

export function isRskRatings(arg: any): arg is models.RskRatings {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // numScores?: number
    ( typeof arg.numScores === 'undefined' || typeof arg.numScores === 'number' ) &&
    // scoreAvg?: number
    ( typeof arg.scoreAvg === 'undefined' || typeof arg.scoreAvg === 'number' ) &&
    // scores?: { [key: string]: number }
    ( typeof arg.scores === 'undefined' || typeof arg.scores === 'number' ) &&

  true
  );
  }

export function isRskRecipientStats(arg: any): arg is models.RskRecipientStats {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // donatedAmountUsd?: number
    ( typeof arg.donatedAmountUsd === 'undefined' || typeof arg.donatedAmountUsd === 'number' ) &&
    // donationRecipient?: string
    ( typeof arg.donationRecipient === 'undefined' || typeof arg.donationRecipient === 'string' ) &&
    // pointsSpent?: number
    ( typeof arg.pointsSpent === 'undefined' || typeof arg.pointsSpent === 'number' ) &&

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

export function isRskRoadmap(arg: any): arg is models.RskRoadmap {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // markdown?: string
    ( typeof arg.markdown === 'undefined' || typeof arg.markdown === 'string' ) &&

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
    // stats?: { [key: string]: RskSearchStats }
    ( typeof arg.stats === 'undefined' || isRskSearchStats(arg.stats) ) &&

  true
  );
  }

export function isRskSearchStats(arg: any): arg is models.RskSearchStats {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // labels?: string[]
    ( typeof arg.labels === 'undefined' || (Array.isArray(arg.labels) && arg.labels.every((item: unknown) => typeof item === 'string')) ) &&
    // values?: number[]
    ( typeof arg.values === 'undefined' || (Array.isArray(arg.values) && arg.values.every((item: unknown) => typeof item === 'number')) ) &&

  true
  );
  }

export function isRskSearchTermPredictions(arg: any): arg is models.RskSearchTermPredictions {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // predictions?: RskPrediction[]
    ( typeof arg.predictions === 'undefined' || (Array.isArray(arg.predictions) && arg.predictions.every((item: unknown) => isRskPrediction(item))) ) &&
    // prefix?: string
    ( typeof arg.prefix === 'undefined' || typeof arg.prefix === 'string' ) &&

  true
  );
  }

export function isRskShortTranscript(arg: any): arg is models.RskShortTranscript {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // actors?: string[]
    ( typeof arg.actors === 'undefined' || (Array.isArray(arg.actors) && arg.actors.every((item: unknown) => typeof item === 'string')) ) &&
    // audioQuality?: RskAudioQuality
    ( typeof arg.audioQuality === 'undefined' || isRskAudioQuality(arg.audioQuality) ) &&
    // bestof?: boolean
    ( typeof arg.bestof === 'undefined' || typeof arg.bestof === 'boolean' ) &&
    // episode?: number
    ( typeof arg.episode === 'undefined' || typeof arg.episode === 'number' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // incomplete?: boolean
    ( typeof arg.incomplete === 'undefined' || typeof arg.incomplete === 'boolean' ) &&
    // media?: RskMedia
    ( typeof arg.media === 'undefined' || isRskMedia(arg.media) ) &&
    // metadata?: { [key: string]: string }
    ( typeof arg.metadata === 'undefined' || typeof arg.metadata === 'string' ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // numRatingScores?: number
    ( typeof arg.numRatingScores === 'undefined' || typeof arg.numRatingScores === 'number' ) &&
    // offsetAccuracyPcnt?: number
    ( typeof arg.offsetAccuracyPcnt === 'undefined' || typeof arg.offsetAccuracyPcnt === 'number' ) &&
    // publication?: string
    ( typeof arg.publication === 'undefined' || typeof arg.publication === 'string' ) &&
    // publicationType?: RskPublicationType
    ( typeof arg.publicationType === 'undefined' || isRskPublicationType(arg.publicationType) ) &&
    // ratingBreakdown?: { [key: string]: number }
    ( typeof arg.ratingBreakdown === 'undefined' || typeof arg.ratingBreakdown === 'number' ) &&
    // ratingScore?: number
    ( typeof arg.ratingScore === 'undefined' || typeof arg.ratingScore === 'number' ) &&
    // releaseDate?: string
    ( typeof arg.releaseDate === 'undefined' || typeof arg.releaseDate === 'string' ) &&
    // series?: number
    ( typeof arg.series === 'undefined' || typeof arg.series === 'number' ) &&
    // shortId?: string
    ( typeof arg.shortId === 'undefined' || typeof arg.shortId === 'string' ) &&
    // special?: boolean
    ( typeof arg.special === 'undefined' || typeof arg.special === 'boolean' ) &&
    // summary?: string
    ( typeof arg.summary === 'undefined' || typeof arg.summary === 'string' ) &&
    // synopsis?: RskSynopsis[]
    ( typeof arg.synopsis === 'undefined' || (Array.isArray(arg.synopsis) && arg.synopsis.every((item: unknown) => isRskSynopsis(item))) ) &&
    // transcriptAvailable?: boolean
    ( typeof arg.transcriptAvailable === 'undefined' || typeof arg.transcriptAvailable === 'boolean' ) &&
    // triviaAvailable?: boolean
    ( typeof arg.triviaAvailable === 'undefined' || typeof arg.triviaAvailable === 'boolean' ) &&
    // version?: string
    ( typeof arg.version === 'undefined' || typeof arg.version === 'string' ) &&

  true
  );
  }

export function isRskShortTranscriptChange(arg: any): arg is models.RskShortTranscriptChange {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // author?: RskAuthor
    ( typeof arg.author === 'undefined' || isRskAuthor(arg.author) ) &&
    // createdAt?: string
    ( typeof arg.createdAt === 'undefined' || typeof arg.createdAt === 'string' ) &&
    // episodeId?: string
    ( typeof arg.episodeId === 'undefined' || typeof arg.episodeId === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // merged?: boolean
    ( typeof arg.merged === 'undefined' || typeof arg.merged === 'boolean' ) &&
    // pointsAwarded?: number
    ( typeof arg.pointsAwarded === 'undefined' || typeof arg.pointsAwarded === 'number' ) &&
    // state?: RskContributionState
    ( typeof arg.state === 'undefined' || isRskContributionState(arg.state) ) &&
    // transcriptVersion?: string
    ( typeof arg.transcriptVersion === 'undefined' || typeof arg.transcriptVersion === 'string' ) &&

  true
  );
  }

export function isRskSong(arg: any): arg is models.RskSong {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // album?: string
    ( typeof arg.album === 'undefined' || typeof arg.album === 'string' ) &&
    // albumImageUrl?: string
    ( typeof arg.albumImageUrl === 'undefined' || typeof arg.albumImageUrl === 'string' ) &&
    // artist?: string
    ( typeof arg.artist === 'undefined' || typeof arg.artist === 'string' ) &&
    // episodeIds?: string[]
    ( typeof arg.episodeIds === 'undefined' || (Array.isArray(arg.episodeIds) && arg.episodeIds.every((item: unknown) => typeof item === 'string')) ) &&
    // spotifyUri?: string
    ( typeof arg.spotifyUri === 'undefined' || typeof arg.spotifyUri === 'string' ) &&
    // title?: string
    ( typeof arg.title === 'undefined' || typeof arg.title === 'string' ) &&
    // transcribed?: string[]
    ( typeof arg.transcribed === 'undefined' || (Array.isArray(arg.transcribed) && arg.transcribed.every((item: unknown) => typeof item === 'string')) ) &&

  true
  );
  }

export function isRskSongList(arg: any): arg is models.RskSongList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // resultCount?: number
    ( typeof arg.resultCount === 'undefined' || typeof arg.resultCount === 'number' ) &&
    // songs?: RskSong[]
    ( typeof arg.songs === 'undefined' || (Array.isArray(arg.songs) && arg.songs.every((item: unknown) => isRskSong(item))) ) &&

  true
  );
  }

export function isRskSynopsis(arg: any): arg is models.RskSynopsis {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // description?: string
    ( typeof arg.description === 'undefined' || typeof arg.description === 'string' ) &&
    // endPos?: number
    ( typeof arg.endPos === 'undefined' || typeof arg.endPos === 'number' ) &&
    // startPos?: number
    ( typeof arg.startPos === 'undefined' || typeof arg.startPos === 'number' ) &&

  true
  );
  }

export function isRskTag(arg: any): arg is models.RskTag {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // timestamp?: string
    ( typeof arg.timestamp === 'undefined' || typeof arg.timestamp === 'string' ) &&

  true
  );
  }

export function isRskTranscript(arg: any): arg is models.RskTranscript {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // actors?: string[]
    ( typeof arg.actors === 'undefined' || (Array.isArray(arg.actors) && arg.actors.every((item: unknown) => typeof item === 'string')) ) &&
    // audioQuality?: RskAudioQuality
    ( typeof arg.audioQuality === 'undefined' || isRskAudioQuality(arg.audioQuality) ) &&
    // bestof?: boolean
    ( typeof arg.bestof === 'undefined' || typeof arg.bestof === 'boolean' ) &&
    // completionReports?: string[]
    ( typeof arg.completionReports === 'undefined' || (Array.isArray(arg.completionReports) && arg.completionReports.every((item: unknown) => typeof item === 'string')) ) &&
    // contributors?: string[]
    ( typeof arg.contributors === 'undefined' || (Array.isArray(arg.contributors) && arg.contributors.every((item: unknown) => typeof item === 'string')) ) &&
    // episode?: number
    ( typeof arg.episode === 'undefined' || typeof arg.episode === 'number' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // incomplete?: boolean
    ( typeof arg.incomplete === 'undefined' || typeof arg.incomplete === 'boolean' ) &&
    // locked?: boolean
    ( typeof arg.locked === 'undefined' || typeof arg.locked === 'boolean' ) &&
    // media?: RskMedia
    ( typeof arg.media === 'undefined' || isRskMedia(arg.media) ) &&
    // metadata?: { [key: string]: string }
    ( typeof arg.metadata === 'undefined' || typeof arg.metadata === 'string' ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // offsetAccuracyPcnt?: number
    ( typeof arg.offsetAccuracyPcnt === 'undefined' || typeof arg.offsetAccuracyPcnt === 'number' ) &&
    // publication?: string
    ( typeof arg.publication === 'undefined' || typeof arg.publication === 'string' ) &&
    // publicationType?: RskPublicationType
    ( typeof arg.publicationType === 'undefined' || isRskPublicationType(arg.publicationType) ) &&
    // ratings?: RskRatings
    ( typeof arg.ratings === 'undefined' || isRskRatings(arg.ratings) ) &&
    // rawTranscript?: string
    ( typeof arg.rawTranscript === 'undefined' || typeof arg.rawTranscript === 'string' ) &&
    // releaseDate?: string
    ( typeof arg.releaseDate === 'undefined' || typeof arg.releaseDate === 'string' ) &&
    // series?: number
    ( typeof arg.series === 'undefined' || typeof arg.series === 'number' ) &&
    // shortId?: string
    ( typeof arg.shortId === 'undefined' || typeof arg.shortId === 'string' ) &&
    // special?: boolean
    ( typeof arg.special === 'undefined' || typeof arg.special === 'boolean' ) &&
    // summary?: string
    ( typeof arg.summary === 'undefined' || typeof arg.summary === 'string' ) &&
    // synopses?: RskSynopsis[]
    ( typeof arg.synopses === 'undefined' || (Array.isArray(arg.synopses) && arg.synopses.every((item: unknown) => isRskSynopsis(item))) ) &&
    // tags?: RskTag[]
    ( typeof arg.tags === 'undefined' || (Array.isArray(arg.tags) && arg.tags.every((item: unknown) => isRskTag(item))) ) &&
    // transcript?: RskDialog[]
    ( typeof arg.transcript === 'undefined' || (Array.isArray(arg.transcript) && arg.transcript.every((item: unknown) => isRskDialog(item))) ) &&
    // trivia?: RskTrivia[]
    ( typeof arg.trivia === 'undefined' || (Array.isArray(arg.trivia) && arg.trivia.every((item: unknown) => isRskTrivia(item))) ) &&
    // version?: string
    ( typeof arg.version === 'undefined' || typeof arg.version === 'string' ) &&

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
    // episodeId?: string
    ( typeof arg.episodeId === 'undefined' || typeof arg.episodeId === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // merged?: boolean
    ( typeof arg.merged === 'undefined' || typeof arg.merged === 'boolean' ) &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // pointsAwarded?: number
    ( typeof arg.pointsAwarded === 'undefined' || typeof arg.pointsAwarded === 'number' ) &&
    // state?: RskContributionState
    ( typeof arg.state === 'undefined' || isRskContributionState(arg.state) ) &&
    // summary?: string
    ( typeof arg.summary === 'undefined' || typeof arg.summary === 'string' ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&
    // transcriptVersion?: string
    ( typeof arg.transcriptVersion === 'undefined' || typeof arg.transcriptVersion === 'string' ) &&

  true
  );
  }

export function isRskTranscriptChangeDiff(arg: any): arg is models.RskTranscriptChangeDiff {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // diffs?: string[]
    ( typeof arg.diffs === 'undefined' || (Array.isArray(arg.diffs) && arg.diffs.every((item: unknown) => typeof item === 'string')) ) &&

  true
  );
  }

export function isRskTranscriptChangeList(arg: any): arg is models.RskTranscriptChangeList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // changes?: RskShortTranscriptChange[]
    ( typeof arg.changes === 'undefined' || (Array.isArray(arg.changes) && arg.changes.every((item: unknown) => isRskShortTranscriptChange(item))) ) &&

  true
  );
  }

export function isRskTranscriptChunkList(arg: any): arg is models.RskTranscriptChunkList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // chunks?: RskChunk[]
    ( typeof arg.chunks === 'undefined' || (Array.isArray(arg.chunks) && arg.chunks.every((item: unknown) => isRskChunk(item))) ) &&

  true
  );
  }

export function isRskTranscriptDialog(arg: any): arg is models.RskTranscriptDialog {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // dialog?: RskDialog[]
    ( typeof arg.dialog === 'undefined' || (Array.isArray(arg.dialog) && arg.dialog.every((item: unknown) => isRskDialog(item))) ) &&
    // maxDialogPosition?: number
    ( typeof arg.maxDialogPosition === 'undefined' || typeof arg.maxDialogPosition === 'number' ) &&
    // transcriptMeta?: RskShortTranscript
    ( typeof arg.transcriptMeta === 'undefined' || isRskShortTranscript(arg.transcriptMeta) ) &&

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

export function isRskTrivia(arg: any): arg is models.RskTrivia {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // description?: string
    ( typeof arg.description === 'undefined' || typeof arg.description === 'string' ) &&
    // endPos?: number
    ( typeof arg.endPos === 'undefined' || typeof arg.endPos === 'number' ) &&
    // startPos?: number
    ( typeof arg.startPos === 'undefined' || typeof arg.startPos === 'number' ) &&

  true
  );
  }

export function isRskTscriptImport(arg: any): arg is models.RskTscriptImport {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // completedAt?: string
    ( typeof arg.completedAt === 'undefined' || typeof arg.completedAt === 'string' ) &&
    // createdAt?: string
    ( typeof arg.createdAt === 'undefined' || typeof arg.createdAt === 'string' ) &&
    // epid?: string
    ( typeof arg.epid === 'undefined' || typeof arg.epid === 'string' ) &&
    // epname?: string
    ( typeof arg.epname === 'undefined' || typeof arg.epname === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // log?: RskTscriptImportLog[]
    ( typeof arg.log === 'undefined' || (Array.isArray(arg.log) && arg.log.every((item: unknown) => isRskTscriptImportLog(item))) ) &&
    // mp3Uri?: string
    ( typeof arg.mp3Uri === 'undefined' || typeof arg.mp3Uri === 'string' ) &&

  true
  );
  }

export function isRskTscriptImportList(arg: any): arg is models.RskTscriptImportList {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // imports?: RskTscriptImport[]
    ( typeof arg.imports === 'undefined' || (Array.isArray(arg.imports) && arg.imports.every((item: unknown) => isRskTscriptImport(item))) ) &&

  true
  );
  }

export function isRskTscriptImportLog(arg: any): arg is models.RskTscriptImportLog {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // msg?: string
    ( typeof arg.msg === 'undefined' || typeof arg.msg === 'string' ) &&
    // stage?: string
    ( typeof arg.stage === 'undefined' || typeof arg.stage === 'string' ) &&

  true
  );
  }

export function isTranscriptServiceBulkSetTranscriptRatingScoreBody(arg: any): arg is models.TranscriptServiceBulkSetTranscriptRatingScoreBody {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // oauthSource?: string
    ( typeof arg.oauthSource === 'undefined' || typeof arg.oauthSource === 'string' ) &&
    // scores?: { [key: string]: number }
    ( typeof arg.scores === 'undefined' || typeof arg.scores === 'number' ) &&

  true
  );
  }

export function isTranscriptServiceBulkSetTranscriptTagsBody(arg: any): arg is models.TranscriptServiceBulkSetTranscriptTagsBody {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // tags?: RskTag[]
    ( typeof arg.tags === 'undefined' || (Array.isArray(arg.tags) && arg.tags.every((item: unknown) => isRskTag(item))) ) &&

  true
  );
  }

export function isTranscriptServiceCreateChunkContributionBody(arg: any): arg is models.TranscriptServiceCreateChunkContributionBody {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }

export function isTranscriptServiceCreateTranscriptChangeBody(arg: any): arg is models.TranscriptServiceCreateTranscriptChangeBody {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // summary?: string
    ( typeof arg.summary === 'undefined' || typeof arg.summary === 'string' ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&
    // transcriptVersion?: string
    ( typeof arg.transcriptVersion === 'undefined' || typeof arg.transcriptVersion === 'string' ) &&

  true
  );
  }

export function isTranscriptServiceRequestChunkContributionStateBody(arg: any): arg is models.TranscriptServiceRequestChunkContributionStateBody {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // comment?: string
    ( typeof arg.comment === 'undefined' || typeof arg.comment === 'string' ) &&
    // requestState?: RskContributionState
    ( typeof arg.requestState === 'undefined' || isRskContributionState(arg.requestState) ) &&

  true
  );
  }

export function isTranscriptServiceRequestTranscriptChangeStateBody(arg: any): arg is models.TranscriptServiceRequestTranscriptChangeStateBody {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // pointsOnApprove?: number
    ( typeof arg.pointsOnApprove === 'undefined' || typeof arg.pointsOnApprove === 'number' ) &&
    // state?: RskContributionState
    ( typeof arg.state === 'undefined' || isRskContributionState(arg.state) ) &&

  true
  );
  }

export function isTranscriptServiceSetTranscriptRatingScoreBody(arg: any): arg is models.TranscriptServiceSetTranscriptRatingScoreBody {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // delete?: boolean
    ( typeof arg.delete === 'undefined' || typeof arg.delete === 'boolean' ) &&
    // score?: number
    ( typeof arg.score === 'undefined' || typeof arg.score === 'number' ) &&

  true
  );
  }

export function isTranscriptServiceUpdateChunkContributionBody(arg: any): arg is models.TranscriptServiceUpdateChunkContributionBody {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // state?: RskContributionState
    ( typeof arg.state === 'undefined' || isRskContributionState(arg.state) ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }

export function isTranscriptServiceUpdateTranscriptChangeBody(arg: any): arg is models.TranscriptServiceUpdateTranscriptChangeBody {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // name?: string
    ( typeof arg.name === 'undefined' || typeof arg.name === 'string' ) &&
    // pointsOnApprove?: number
    ( typeof arg.pointsOnApprove === 'undefined' || typeof arg.pointsOnApprove === 'number' ) &&
    // state?: RskContributionState
    ( typeof arg.state === 'undefined' || isRskContributionState(arg.state) ) &&
    // summary?: string
    ( typeof arg.summary === 'undefined' || typeof arg.summary === 'string' ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }


