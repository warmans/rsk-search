/* tslint:disable */
import {
  RskCurrentRadioEpisode,
} from '.';

export interface RskPutRadioStateRequest {
  currentEpisode?: RskCurrentRadioEpisode;
  currentTimestampMs?: number;
}
