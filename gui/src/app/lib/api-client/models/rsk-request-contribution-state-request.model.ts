/* tslint:disable */
import {
  RskContributionState,
} from '.';

export interface RskRequestContributionStateRequest {
  comment?: string;
  contributionId?: string;
  requestState?: RskContributionState;
}
