/* tslint:disable */
import {
  RsksearchAuthor,
  RsksearchContributionState,
} from '.';

export interface RsksearchTscriptContribution {
  author?: RsksearchAuthor;
  chunkId?: string;
  createdAt?: string;
  id?: string;
  state?: RsksearchContributionState;
  transcript?: string;
  tscriptId?: string;
}
