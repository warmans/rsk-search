/* tslint:disable */
import {
  RskAuthor,
  RskContributionState,
} from '.';

export interface RskChunkContribution {
  author?: RskAuthor;
  chunkId?: string;
  id?: string;
  state?: RskContributionState;
  transcript?: string;
}
