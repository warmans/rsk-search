/* tslint:disable */
import {
  RskContributionState,
} from '.';

export interface RskUpdateContributionRequest {
  contributionId?: string;
  state?: RskContributionState;
  transcript?: string;
}
