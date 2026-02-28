/* tslint:disable */
import {
  RskContributionState,
} from '.';

export interface TranscriptServiceUpdateTranscriptChangeBody {
  name?: string;
  pointsOnApprove?: number;
  releaseDate?: string;
  state?: RskContributionState;
  summary?: string;
  transcript?: string;
}
