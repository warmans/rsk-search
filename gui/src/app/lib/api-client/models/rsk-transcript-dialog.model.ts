/* tslint:disable */
import {
  RskDialog,
  RskShortTranscript,
} from '.';

export interface RskTranscriptDialog {
  dialog?: RskDialog[];
  maxDialogPosition?: number;
  transcriptMeta?: RskShortTranscript;
}
