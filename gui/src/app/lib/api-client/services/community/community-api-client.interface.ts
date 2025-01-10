/* tslint:disable */

import { Observable } from 'rxjs';
import { HttpOptions } from './';
import * as models from '../../models';

export interface CommunityAPIClientInterface {

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listArchive(
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskArchiveList>;

  /**
   * Response generated for [ 200 ] HTTP response code.
   */
  listCommunityProjects(
    args: {
      filter?: string,
      sortField?: string,
      sortDirection?: string,
      page?: number,
      pageSize?: number,
    },
    requestHttpOptions?: HttpOptions
  ): Observable<models.RskCommunityProjectList>;

}
