/* tslint:disable */
import {
  RskContributionState,
} from '.';

export interface RskUpdateChunkContributionRequest {
  chunkId?: string;
  contributionId?: string;
  state?: RskContributionState;
  transcript?: string;
}
