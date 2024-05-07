import {Section} from "../component/transcript/transcript.component";
import {RskDialog} from "../../../lib/api-client/models";

// e.g. pos-1-2
export function parseSection(fragment: string, dialog: RskDialog[]): Section {
  if (!(fragment || '').startsWith("pos-")) {
    return null
  }
  let sections: string[] = fragment.split("-");
  switch (sections.length) {
    case 2:
      // e.g. pos-1
      const pos = parseInt(sections[1]);
      const posLine = dialog.length > pos - 1 ? dialog[pos - 1] : null;
      return {
        startPos: pos,
        startTimestampMs: posLine ? dialog[pos - 1].offsetMs : 0,
        endPos: pos + 1,
        endTimestampMs: posLine ? dialog[pos - 1].offsetMs + dialog[pos - 1].durationMs : 0,
      }
    case 3:
      const startPos = parseInt(sections[1]);
      const startLine = dialog.length > startPos - 1 ? dialog[startPos - 1] : null;
      const endPos = parseInt(sections[2]);
      const endLine = dialog.length > endPos - 1 ? dialog[endPos - 1] : null;
      // e.g. pos-1-2
      return {
        startPos: startPos,
        startTimestampMs: startLine ? startLine.offsetMs : 0,
        endPos: endPos,
        endTimestampMs: endLine ? endLine.offsetMs + endLine.durationMs : 0,
      }
    default:
      return null
  }
}
