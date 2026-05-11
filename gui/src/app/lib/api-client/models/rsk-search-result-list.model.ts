/* tslint:disable */
import {
  RskSearchResult,
  RskSearchStats,
} from '.';

export interface RskSearchResultList {
  resultCount?: number;
  results?: RskSearchResult[];
  stats?: { [key: string]: RskSearchStats };
}
