/* tslint:disable */
import {
  RsksearchContributionState,
} from '.';

export interface RsksearchChunkContribution {
  authorId?: string;
  chunkId?: string;
  id?: string;
  state?: RsksearchContributionState;
  transcript?: string;
}
