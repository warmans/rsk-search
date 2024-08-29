/* tslint:disable */
import {
  RskCurrentEpisode,
} from '.';

export interface RskPutStateRequest {
  currentEpisode?: RskCurrentEpisode;
  currentTimestampMs?: number;
}
