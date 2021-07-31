/* tslint:disable */
import {
  RskDialog,
  RskSynopsis,
  RskTag,
} from '.';

export interface RskEpisode {
  contributors?: string[];
  episode?: number;
  id?: string;
  incomplete?: boolean;
  metadata?: { [key: string]: string };
  publication?: string;
  releaseDate?: string;
  series?: number;
  synopses?: RskSynopsis[];
  tags?: RskTag[];
  transcript?: RskDialog[];
}
