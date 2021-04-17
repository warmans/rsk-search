/* tslint:disable */
import {
  RskChunkStates,
} from '.';

export interface RskTscriptStats {
  chunkContributions?: { [key: string]: RskChunkStates };
  episode?: number;
  id?: string;
  numApprovedContributions?: number;
  numChunks?: number;
  numContributions?: number;
  numPendingContributions?: number;
  numRejectedContributions?: number;
  numRequestApprovalContributions?: number;
  publication?: string;
  series?: number;
}
