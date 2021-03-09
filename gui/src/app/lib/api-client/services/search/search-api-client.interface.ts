/* tslint:disable */

import { Observable } from 'rxjs';
import { HttpOptions } from './';
import * as models from '../../models';

export interface SearchAPIClientInterface {

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceGetEpisode(
    args: {
      id: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchEpisode>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  searchServiceListFieldValues(
    args: {
      field?: string,
      prefix?: string,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RsksearchFieldValueList>;

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
