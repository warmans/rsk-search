/* tslint:disable */
import {
  RskDialog,
  RskSynopsis,
  RskTrivia,
} from '.';

export interface RskTranscript {
  actors?: string[];
  audioUri?: string;
  contributors?: string[];
  episode?: number;
  id?: string;
  incomplete?: boolean;
  metadata?: { [key: string]: string };
  offsetAccuracyPcnt?: number;
  publication?: string;
  rawTranscript?: string;
  releaseDate?: string;
  series?: number;
  shortId?: string;
  synopses?: RskSynopsis[];
  transcript?: RskDialog[];
  trivia?: RskTrivia[];
}
