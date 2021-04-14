/* tslint:disable */
import {
  RsksearchAuthor,
  RsksearchContributionState,
} from '.';

export interface RsksearchChunkContribution {
  author?: RsksearchAuthor;
  chunkId?: string;
  id?: string;
  state?: RsksearchContributionState;
  transcript?: string;
}
