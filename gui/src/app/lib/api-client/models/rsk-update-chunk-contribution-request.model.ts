/* tslint:disable */
import {
  RskContributionState,
} from '.';

export interface RskUpdateChunkContributionRequest {
  contributionId?: string;
  state?: RskContributionState;
  transcript?: string;
}
