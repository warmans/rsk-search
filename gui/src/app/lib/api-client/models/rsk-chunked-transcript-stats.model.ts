/* tslint:disable */
import {
  RskChunkStates,
} from '.';

export interface RskChunkedTranscriptStats {
  chunkContributions?: { [key: string]: RskChunkStates };
  episode?: number;
  id?: string;
  name?: string;
  numApprovedContributions?: number;
  numChunks?: number;
  numContributions?: number;
  numPendingContributions?: number;
  numRejectedContributions?: number;
  numRequestApprovalContributions?: number;
  publication?: string;
  series?: number;
}
