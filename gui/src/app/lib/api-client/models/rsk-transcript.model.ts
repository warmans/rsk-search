/* tslint:disable */
import {
  RskAudioQuality,
  RskDialog,
  RskMedia,
  RskMediaType,
  RskPublicationType,
  RskRatings,
  RskSynopsis,
  RskTrivia,
} from '.';

export interface RskTranscript {
  actors?: string[];
  audioQuality?: RskAudioQuality;
  audioUri?: string;
  bestof?: boolean;
  contributors?: string[];
  episode?: number;
  id?: string;
  incomplete?: boolean;
  locked?: boolean;
  media?: RskMedia;
  mediaType?: RskMediaType;
  metadata?: { [key: string]: string };
  name?: string;
  offsetAccuracyPcnt?: number;
  publication?: string;
  publicationType?: RskPublicationType;
  ratings?: RskRatings;
  rawTranscript?: string;
  releaseDate?: string;
  series?: number;
  shortId?: string;
  special?: boolean;
  summary?: string;
  synopses?: RskSynopsis[];
  transcript?: RskDialog[];
  trivia?: RskTrivia[];
  version?: string;
}
