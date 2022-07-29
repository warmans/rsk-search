/* tslint:disable */
import {
  RskAuthor,
  RskContributionState,
} from '.';

export interface RskTranscriptChange {
  author?: RskAuthor;
  createdAt?: string;
  episodeId?: string;
  id?: string;
  merged?: boolean;
  pointsAwarded?: number;
  state?: RskContributionState;
  summary?: string;
  transcript?: string;
  transcriptVersion?: string;
}
