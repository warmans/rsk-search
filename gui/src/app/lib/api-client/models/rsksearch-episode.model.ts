/* tslint:disable */
import {
  RsksearchDialog,
} from '.';

export interface RsksearchEpisode {
  episode?: number;
  id?: string;
  metadata?: { [key: string]: string };
  publication?: string;
  series?: number;
  tags?: { [key: string]: string };
  transcript?: RsksearchDialog[];
}
