/* tslint:disable */
import {
  RsksearchChunkStates,
} from '.';

export interface RsksearchTscriptStats {
  chunkContributions?: { [key: string]: RsksearchChunkStates };
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
