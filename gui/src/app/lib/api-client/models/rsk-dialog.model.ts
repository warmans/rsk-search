/* tslint:disable */
import {
  DialogType,
} from '.';

export interface RskDialog {
  actor?: string;
  content?: string;
  durationMs?: number;
  id?: string;
  isMatchedRow?: boolean;
  metadata?: { [key: string]: string };
  notable?: boolean;
  offsetDistance?: number;
  offsetInferred?: boolean;
  offsetMs?: number;
  offsetSec?: string;
  placeholder?: boolean;
  pos?: number;
  type?: DialogType;
}
