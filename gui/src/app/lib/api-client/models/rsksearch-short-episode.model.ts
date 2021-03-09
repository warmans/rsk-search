/* tslint:disable */
import {
  RsksearchTag,
} from '.';

export interface RsksearchShortEpisode {
  episode?: number;
  id?: string;
  metadata?: { [key: string]: string };
  publication?: string;
  series?: number;
  tags?: RsksearchTag[];
}
