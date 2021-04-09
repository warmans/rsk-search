/* tslint:disable */
import {
  RewardKind,
} from '.';

export interface RsksearchReward {
  criteria?: string;
  id?: string;
  kind?: RewardKind;
  name?: string;
  value?: number;
  valueCurrency?: string;
}
