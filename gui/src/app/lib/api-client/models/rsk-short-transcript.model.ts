/* tslint:disable */
import {
  RskAudioQuality,
  RskMedia,
  RskPublicationType,
  RskSynopsis,
} from '.';

export interface RskShortTranscript {
  actors?: string[];
  audioQuality?: RskAudioQuality;
  bestof?: boolean;
  episode?: number;
  id?: string;
  incomplete?: boolean;
  media?: RskMedia;
  metadata?: { [key: string]: string };
  name?: string;
  numRatingScores?: number;
  offsetAccuracyPcnt?: number;
  publication?: string;
  publicationType?: RskPublicationType;
  ratingScore?: number;
  releaseDate?: string;
  series?: number;
  shortId?: string;
  special?: boolean;
  summary?: string;
  synopsis?: RskSynopsis[];
  transcriptAvailable?: boolean;
  triviaAvailable?: boolean;
  version?: string;
}
