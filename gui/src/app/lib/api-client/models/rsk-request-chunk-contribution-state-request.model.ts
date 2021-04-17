/* tslint:disable */
import {
  RskContributionState,
} from '.';

export interface RskRequestChunkContributionStateRequest {
  chunkId?: string;
  comment?: string;
  contributionId?: string;
  requestState?: RskContributionState;
}
