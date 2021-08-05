/* tslint:disable */
import {
  RskContributionState,
} from '.';

export interface RskUpdateTranscriptChangeRequest {
  id?: string;
  state?: RskContributionState;
  transcript?: string;
}
