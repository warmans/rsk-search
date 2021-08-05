/* tslint:disable */
import {
  RskDialog,
  RskSynopsis,
  RskTag,
} from '.';

export interface RskTranscript {
  audioUri?: string;
  contributors?: string[];
  episode?: number;
  id?: string;
  incomplete?: boolean;
  metadata?: { [key: string]: string };
  publication?: string;
  rawTranscript?: string;
  releaseDate?: string;
  series?: number;
  shortId?: string;
  synopses?: RskSynopsis[];
  tags?: RskTag[];
  transcript?: RskDialog[];
}
