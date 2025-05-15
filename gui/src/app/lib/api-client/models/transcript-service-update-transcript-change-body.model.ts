/* tslint:disable */
import {
  RskContributionState,
} from '.';

export interface TranscriptServiceUpdateTranscriptChangeBody {
  name?: string;
  pointsOnApprove?: number;
  state?: RskContributionState;
  summary?: string;
  transcript?: string;
}
