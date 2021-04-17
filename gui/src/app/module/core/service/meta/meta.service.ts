import { Injectable } from '@angular/core';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { Observable, of } from 'rxjs';
import { FieldMetaKind, RskSearchMetadata } from '../../../../lib/api-client/models';
import { first } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class MetaService {

  private cache: RskSearchMetadata;

  constructor(private apiClient: SearchAPIClient) {
    this.apiClient.getSearchMetadata().pipe(first()).subscribe((v: RskSearchMetadata) => {
      this.cache = v;
    });
  }

  getMeta(): Observable<RskSearchMetadata> {
    if (this.cache === undefined) {
      return this.apiClient.getSearchMetadata();
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
