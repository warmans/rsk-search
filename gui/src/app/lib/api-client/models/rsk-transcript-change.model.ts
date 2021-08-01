/* tslint:disable */
import {
  RskAuthor,
  RskContributionState,
} from '.';

export interface RskTranscriptChange {
  author?: RskAuthor;
  createdAt?: string;
  diff?: string;
  episodeId?: string;
  id?: string;
  state?: RskContributionState;
  transcript?: string;
}
