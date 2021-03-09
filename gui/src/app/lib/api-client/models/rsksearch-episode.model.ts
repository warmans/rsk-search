/* tslint:disable */
import {
  RsksearchDialog,
  RsksearchTag,
} from '.';

export interface RsksearchEpisode {
  episode?: number;
  id?: string;
  metadata?: { [key: string]: string };
  publication?: string;
  series?: number;
  tags?: RsksearchTag[];
  transcript?: RsksearchDialog[];
}
