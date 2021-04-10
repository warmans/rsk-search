/* tslint:disable */
import {
  RsksearchDialog,
  RsksearchSynopsis,
  RsksearchTag,
} from '.';

export interface RsksearchEpisode {
  contributors?: string[];
  episode?: number;
  id?: string;
  metadata?: { [key: string]: string };
  publication?: string;
  releaseDate?: string;
  series?: number;
  synopses?: RsksearchSynopsis[];
  tags?: RsksearchTag[];
  transcript?: RsksearchDialog[];
}
