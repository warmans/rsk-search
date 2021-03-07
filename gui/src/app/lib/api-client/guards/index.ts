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

export function isRsksearchDialog(arg: any): arg is models.RsksearchDialog {
  return (
  arg != null &&
  typeof arg === 'object' &&
    // actor?: string
    ( typeof arg.actor === 'undefined' || typeof arg.actor === 'string' ) &&
    // content?: string
    ( typeof arg.content === 'undefined' || typeof arg.content === 'string' ) &&
    // id?: string
    ( typeof arg.id === 'undefined' || typeof arg.id === 'string' ) &&
    // isMatchedRow?: boolean
    ( typeof arg.isMatchedRow === 'undefined' || typeof arg.isMatchedRow === 'boolean' ) &&
    // metadata?: { [key: string]: string }
    ( typeof arg.metadata === 'undefined' || typeof arg.metadata === 'string' ) &&
    // pos?: string
    ( typeof arg.pos === 'undefined' || typeof arg.pos === 'string' ) &&
    // tags?: { [key: string]: string }
    ( typeof arg.tags === 'undefined' || typeof arg.tags === 'string' ) &&
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
    // series?: number
    ( typeof arg.series === 'undefined' || typeof arg.series === 'number' ) &&
    // tags?: { [key: string]: string }
    ( typeof arg.tags === 'undefined' || typeof arg.tags === 'string' ) &&
    // transcript?: RsksearchDialog[]
    ( typeof arg.transcript === 'undefined' || (Array.isArray(arg.transcript) && arg.transcript.every((item: unknown) => isRsksearchDialog(item))) ) &&

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

export function isRsksearchShortEpisode(arg: any): arg is models.RsksearchShortEpisode {
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
    // series?: number
    ( typeof arg.series === 'undefined' || typeof arg.series === 'number' ) &&

  true
  );
  }


