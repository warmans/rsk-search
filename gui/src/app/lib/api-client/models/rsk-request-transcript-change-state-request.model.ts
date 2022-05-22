/* tslint:disable */
import {
  RskContributionState,
} from '.';

export interface RskRequestTranscriptChangeStateRequest {
  id?: string;
  pointsOnApprove?: number;
  state?: RskContributionState;
}
