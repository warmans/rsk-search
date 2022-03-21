/* tslint:disable */
import {
  RskAuthor,
  RskRank,
} from '.';

export interface RskAuthorRank {
  approvedChanges?: number;
  approvedChunks?: number;
  author?: RskAuthor;
  currentRank?: RskRank;
  nextRank?: RskRank;
  points?: number;
  rewardValueUsd?: number;
}
