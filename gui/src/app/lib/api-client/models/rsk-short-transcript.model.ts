/* tslint:disable */
import {
  RskSynopsis,
} from '.';

export interface RskShortTranscript {
  actors?: string[];
  episode?: number;
  id?: string;
  incomplete?: boolean;
  publication?: string;
  releaseDate?: string;
  series?: number;
  summary?: string;
  synopsis?: RskSynopsis[];
  transcriptAvailable?: boolean;
  triviaAvailable?: boolean;
}
