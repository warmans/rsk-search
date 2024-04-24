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
  offsetDistance?: string;
  offsetInferred?: boolean;
  offsetMs?: string;
  offsetSec?: string;
  pos?: number;
  type?: DialogType;
}
