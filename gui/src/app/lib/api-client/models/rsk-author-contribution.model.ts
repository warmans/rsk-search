/* tslint:disable */
import {
  AuthorContributionType,
  RskAuthor,
} from '.';

export interface RskAuthorContribution {
  author?: RskAuthor;
  contributionType?: AuthorContributionType;
  createdAt?: string;
  episodeId?: string;
  id?: string;
  points?: number;
}
