import {Section} from "../component/transcript/transcript.component";

// e.g. pos-1-2
export function parseSection(fragment: string): Section {
  if (!fragment.startsWith("pos-")) {
    return null
  }
  let sections: string[] = fragment.split("-");
  switch (sections.length) {
    case 2:
      // e.g. pos-1
      return {
        startPos: parseInt(sections[1]),
      }
    case 3:
      // e.g. pos-1-2
      return {
        startPos: parseInt(sections[1]),
        endPos: parseInt(sections[2]),
      }
    default:
      return null
  }
}
