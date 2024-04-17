import {DialogType, RskDialog, RskSynopsis, RskTrivia} from '../../../lib/api-client/models';

export function lineHasActorPrefix(line: string): boolean {
  return line.length > 1 && line.indexOf(':') > -1;
}

export function isOffsetLine(line: string): boolean {
  return getOffsetValueFromLine(line) > -1;
}

export function isSynopsisLine(line: string): boolean {
  return isStartSynopsisLine(line) || isEndSynopsisLine(line);
}

export function getSynopsis(line: string): string {
  const match = line.match(/^#SYN:\s(.+)/);
  return match?.length == 2 ? match[1] : '';
}

export function isStartSynopsisLine(line: string): boolean {
  return !!line.match(/^#SYN:.+/g);
}

export function isEndSynopsisLine(line: string): boolean {
  return !!line.match(/^#[/]SYN.*/g);
}

export function isTriviaLine(line: string): boolean {
  return isStartTriviaLine(line) || isEndTriviaLine(line);
}

// synopsis or trivia text block.
export function isMetadataBlockText(line: string): boolean {
  return !!line.match(/^#.*/g);
}

export function getTrivia(line: string): string {
  const match = line.match(/^#TRIVIA:(.+)/);
  return match?.length == 2 ? match[1] : '';
}

export function isStartTriviaLine(line: string): boolean {
  return !!line.match(/^#TRIVIA:.*/g);
}

export function isEndTriviaLine(line: string): boolean {
  return !!line.match(/^#[/]TRIVIA.*/g);
}

export function getOffsetValueFromLine(line: string): number {
  const match = line.match(/^#OFFSET:\s([0-9\.]+)/);
  return match?.length == 2 ? parseFloat(match[1]) : -1;
}

export function getFirstOffset(transcript: string): number {
  for (let line of transcript.split('\n')) {
    let offset = getOffsetValueFromLine(line);
    if (offset > -1) {
      return offset;
    }
  }
  return -1;
}

export function parseTranscript(transcript: string): Tscript {
  let tscript = new Tscript([], [], []);

  if (!transcript) {
    return tscript;
  }

  let currentSynopsis: RskSynopsis;
  let currentTrivia: RskSynopsis;

  let pos = 1;
  transcript.split('\n').forEach((line) => {
    line = line.trim();
    let notable: boolean = false;
    if (line === '') {
      return;
    }
    if (line[0] === '!') {
      line = line.slice(1);
      notable = true;
    }
    if (isOffsetLine(line)) {
      return;
    }
    if (isStartSynopsisLine(line)) {
      currentSynopsis = {description: getSynopsis(line), startPos: pos};
      return;
    }
    if (isEndSynopsisLine(line) && currentSynopsis) {
      currentSynopsis.endPos = pos;
      tscript.synopses.push(currentSynopsis);
      currentSynopsis = undefined;
      return;
    }
    if (isStartTriviaLine(line)) {
      currentTrivia = {description: getTrivia(line), startPos: pos};
      return;
    }
    if (isEndTriviaLine(line) && currentTrivia) {
      currentTrivia.endPos = pos;
      tscript.trivia.push(currentTrivia);
      currentTrivia = undefined;
      return;
    }
    // further lines prefix with a # are considered more trivia lines
    if (currentTrivia && isMetadataBlockText(line)) {
      const lineWithoutPrefix = line.replace(/^#/g, '');
      currentTrivia.description += `\n${lineWithoutPrefix}`;
      return;
    }

    const parts = line.split(':');
    if (parts.length < 2) {
      tscript.transcript.push({
        type: DialogType.UNKNOWN,
        content: parts.join(':'),
        notable: notable,
        pos: pos,
      });
    } else {
      const actor = parts.shift();
      tscript.transcript.push({
        type: actor.toLowerCase() == 'song' ? DialogType.SONG : DialogType.CHAT,
        actor: actor.toLowerCase() === 'none' ? '' : actor,
        content: parts.join(':'),
        notable: notable,
        pos: pos,
      });
    }

    pos++;
  });

  if (currentTrivia) {
    currentTrivia.endPos = pos;
    tscript.trivia.push(currentTrivia);
  }
  if (currentSynopsis) {
    currentSynopsis.endPos = pos;
    tscript.synopses.push(currentSynopsis);
  }
  return tscript;
}

export class Tscript {
  constructor(public transcript: RskDialog[], public synopses?: RskSynopsis[], public trivia?: RskTrivia[]) {
  }
}
