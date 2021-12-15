import { RskDialog, RskSynopsis, RskTrivia } from '../../../lib/api-client/models';

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

export function getTrivia(line: string): string {
  const match = line.match(/^#TRIVIA:\s(.+)/);
  return match?.length == 2 ? match[1] : '';
}

export function isStartTriviaLine(line: string): boolean {
  return !!line.match(/^#TRIVIA:.+/g);
}

export function isEndTriviaLine(line: string): boolean {
  return !!line.match(/^#[/]TRIVIA.*/g);
}

export function getOffsetValueFromLine(line: string): number {
  const match = line.match(/^#OFFSET:\s([0-9]+)/);
  return match?.length == 2 ? parseInt(match[1], 10) : -1;
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
    if (isOffsetLine(line) || isEndSynopsisLine(line) || isEndTriviaLine(line)) {
      return;
    }
    if (isStartSynopsisLine(line)) {
      // don't bother with positions for now
      tscript.synopses.push({ description: getSynopsis(line) });
      return;
    }
    if (isStartTriviaLine(line)) {
      // don't bother with positions for now
      tscript.trivia.push({ description: getTrivia(line) });
      return;
    }
    const parts = line.split(':');
    if (parts.length < 2) {
      tscript.dialog.push({
        type: 'unknown',
        content: parts.join(':'),
        notable: notable
      });
    } else {
      const actor = parts.shift();
      tscript.dialog.push({
        type: actor.toLowerCase() == 'song' ? 'song' : 'chat',
        actor: actor.toLowerCase() === 'none' ? '' : actor,
        content: parts.join(':'),
        notable: notable
      });
    }
  });

  return tscript;
}

export class Tscript {
  constructor(public dialog: RskDialog[], public synopses: RskSynopsis[], public trivia: RskTrivia[]) {
  }
}
