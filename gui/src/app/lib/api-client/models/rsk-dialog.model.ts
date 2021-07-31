/* tslint:disable */
import {
  RskTag,
} from '.';

export interface RskDialog {
  actor?: string;
  content?: string;
  contentTags?: { [key: string]: RskTag };
  id?: string;
  isMatchedRow?: boolean;
  metadata?: { [key: string]: string };
  notable?: boolean;
  offsetSec?: string;
  pos?: string;
  type?: string;
}
