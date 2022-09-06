/* tslint:disable */
import {
  RskDialog,
  RskSynopsis,
  RskTrivia,
} from '.';

export interface RskTranscript {
  actors?: string[];
  audioUri?: string;
  bestof?: boolean;
  contributors?: string[];
  episode?: number;
  id?: string;
  incomplete?: boolean;
  metadata?: { [key: string]: string };
  name?: string;
  offsetAccuracyPcnt?: number;
  publication?: string;
  rawTranscript?: string;
  releaseDate?: string;
  series?: number;
  shortId?: string;
  special?: boolean;
  synopses?: RskSynopsis[];
  transcript?: RskDialog[];
  trivia?: RskTrivia[];
  version?: string;
}
