/* tslint:disable */
import {
  RskContributionState,
} from '.';

export interface RskUpdateTranscriptChangeRequest {
  id?: string;
  pointsOnApprove?: number;
  state?: RskContributionState;
  summary?: string;
  transcript?: string;
}
