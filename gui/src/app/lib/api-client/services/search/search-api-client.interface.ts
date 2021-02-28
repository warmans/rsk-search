/* tslint:disable */

import { Observable } from 'rxjs';
import { HttpOptions } from './';
import * as models from '../../models';

export interface SearchAPIClientInterface {

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceSearch(
    args: {
      query?: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskSearchResultList>;

}
