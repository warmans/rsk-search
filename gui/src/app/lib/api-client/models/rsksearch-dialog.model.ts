/* tslint:disable */
import {
  RsksearchTag,
} from '.';

export interface RsksearchDialog {
  actor?: string;
  content?: string;
  contentTags?: { [key: string]: RsksearchTag };
  contributor?: string;
  id?: string;
  isMatchedRow?: boolean;
  metadata?: { [key: string]: string };
  notable?: boolean;
  pos?: string;
  type?: string;
}
