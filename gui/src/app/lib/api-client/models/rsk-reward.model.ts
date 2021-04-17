/* tslint:disable */
import {
  RewardKind,
} from '.';

export interface RskReward {
  criteria?: string;
  id?: string;
  kind?: RewardKind;
  name?: string;
  value?: number;
  valueCurrency?: string;
}
