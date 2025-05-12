/* tslint:disable */
import {
  RskAudioQuality,
  RskDialog,
  RskMedia,
  RskPublicationType,
  RskRatings,
  RskSynopsis,
  RskTag,
  RskTrivia,
} from '.';

export interface RskTranscript {
  actors?: string[];
  audioQuality?: RskAudioQuality;
  bestof?: boolean;
  contributors?: string[];
  episode?: number;
  id?: string;
  incomplete?: boolean;
  locked?: boolean;
  media?: RskMedia;
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
  tags?: RskTag[];
  transcript?: RskDialog[];
  trivia?: RskTrivia[];
  version?: string;
}
