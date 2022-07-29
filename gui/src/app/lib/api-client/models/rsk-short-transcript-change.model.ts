/* tslint:disable */
import {
  RskAuthor,
  RskContributionState,
} from '.';

export interface RskShortTranscriptChange {
  author?: RskAuthor;
  createdAt?: string;
  episodeId?: string;
  id?: string;
  merged?: boolean;
  pointsAwarded?: number;
  state?: RskContributionState;
  transcriptVersion?: string;
}
