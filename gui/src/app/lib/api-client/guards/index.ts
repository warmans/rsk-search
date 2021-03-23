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

export function isRequestChunkContributionStateRequestState(arg: any): arg is models.RequestChunkContributionStateRequestState {
  return false
   || arg === models.RequestChunkContributionStateRequestState.STATE_UNDEFINED
   || arg === models.RequestChunkContributionStateRequestState.STATE_REQUEST_APPROVAL
   || arg === models.RequestChunkContributionStateRequestState.STATE_REQUEST_PENDING
  ;
  }

export function isRsksearchChunkContribution(arg: any): arg is models.RsksearchChunkContribution {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // authorId?: string
    ( typeof arg.authorId === 'undefined' || typeof arg.authorId === 'string' ) &&
    // chunkId?: string
    ( typeof arg.chunkId === 'undefined' || typeof arg.chunkId === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // state?: string
    ( typeof arg.state === 'undefined' || typeof arg.state === 'string' ) &&
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
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // isMatchedRow?: boolean
    ( typeof arg.isMatchedRow === 'undefined' || typeof arg.isMatchedRow === 'boolean' ) &&
    // metadata?: { [key: string]: string }
    ( typeof arg.metadata === 'undefined' || typeof arg.metadata === 'string' ) &&
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

export function isRsksearchEpisode(arg: any): arg is models.RsksearchEpisode {
  return (
  arg != null &&
  typeof arg === 'object' &&
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
    // contributionId?: string
    ( typeof arg.contributionId === 'undefined' || typeof arg.contributionId === 'string' ) &&
    // requestState?: RequestChunkContributionStateRequestState
    ( typeof arg.requestState === 'undefined' || isRequestChunkContributionStateRequestState(arg.requestState) ) &&

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
    // state?: string
    ( typeof arg.state === 'undefined' || typeof arg.state === 'string' ) &&

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

export function isRsksearchUpdateChunkContributionRequest(arg: any): arg is models.RsksearchUpdateChunkContributionRequest {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // chunkId?: string
    ( typeof arg.chunkId === 'undefined' || typeof arg.chunkId === 'string' ) &&
    // contributionId?: string
    ( typeof arg.contributionId === 'undefined' || typeof arg.contributionId === 'string' ) &&
    // transcript?: string
    ( typeof arg.transcript === 'undefined' || typeof arg.transcript === 'string' ) &&

  true
  );
  }


