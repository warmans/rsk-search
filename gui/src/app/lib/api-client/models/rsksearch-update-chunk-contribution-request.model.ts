/* tslint:disable */
import {
  RsksearchContributionState,
} from '.';

export interface RsksearchUpdateChunkContributionRequest {
  chunkId?: string;
  contributionId?: string;
  state?: RsksearchContributionState;
  transcript?: string;
}
