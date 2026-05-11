/* tslint:disable */
import {
  RskAuthor,
  RskContributionState,
} from '.';

export interface RskChunkContribution {
  author?: RskAuthor;
  chunkId?: string;
  createdAt?: string;
  id?: string;
  state?: RskContributionState;
  stateComment?: string;
  transcript?: string;
}
