/* tslint:disable */
import {
  RskAuthor,
  RskContributionState,
} from '.';

export interface RskContribution {
  author?: RskAuthor;
  chunkId?: string;
  createdAt?: string;
  id?: string;
  state?: RskContributionState;
  transcript?: string;
  tscriptId?: string;
}
