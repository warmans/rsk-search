/* tslint:disable */
import {
  DialogType,
} from '.';

export interface RskDialog {
  actor?: string;
  content?: string;
  id?: string;
  isMatchedRow?: boolean;
  metadata?: { [key: string]: string };
  notable?: boolean;
  offsetInferred?: boolean;
  offsetSec?: string;
  pos?: number;
  type?: DialogType;
}
