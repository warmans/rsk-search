/* tslint:disable */
import {
  RskTag,
} from '.';

export interface RskDialog {
  actor?: string;
  content?: string;
  contentTags?: { [key: string]: RskTag };
  contributor?: string;
  id?: string;
  isMatchedRow?: boolean;
  metadata?: { [key: string]: string };
  notable?: boolean;
  pos?: string;
  type?: string;
}
