/* tslint:disable */
import {
  RskContributionState,
} from '.';

export interface TranscriptServiceUpdateTranscriptChangeBody {
  pointsOnApprove?: number;
  state?: RskContributionState;
  summary?: string;
  transcript?: string;
}
