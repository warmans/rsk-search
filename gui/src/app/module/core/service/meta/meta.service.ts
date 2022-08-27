import { Injectable } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { Observable, of } from 'rxjs';
import { FieldMetaKind, RskMetadata } from '../../../../lib/api-client/models';
import { first } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class MetaService {

  private cache: RskMetadata;

  constructor(private apiClient: SearchAPIClient) {

  }

  getMeta(): Observable<RskMetadata> {
    if (this.cache === undefined) {
      return this.apiClient.getMetadata();
    }
    return of(this.cache);
  }

  getOperatorsForType(t: FieldMetaKind): string[] {
    switch (t) {
      case FieldMetaKind.IDENTIFIER:
        return ['=', '!='];
      case FieldMetaKind.KEYWORD:
        return ['=', '!='];
      case FieldMetaKind.KEYWORD_LIST:
        return ['=', '!='];
      case FieldMetaKind.TEXT:
        return ['=', '!=', '~='];
      case FieldMetaKind.INT:
        return ['=', '!=', '>', '>=', '<', '<='];
      case FieldMetaKind.FLOAT:
        return ['=', '!=', '>', '>=', '<', '<='];
      case FieldMetaKind.DATE:
        return ['=', '!='];
      case FieldMetaKind.UNKNOWN:
      default:
        return [];
    }
  }
}
