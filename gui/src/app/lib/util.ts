import { RskPrediction, RskWordPosition } from './api-client/models';

export function trimChars(str: string, ch: string): string {
  let start = 0, end = str.length;

  while (start < end && str[start] === ch) {
    ++start;
  }
  while (end > start && str[end - 1] === ch) {
    --end;
  }
  return (start > 0 || end < str.length) ? str.substring(start, end) : str;
}

export function highlightPrediction(pred: RskPrediction): string {
  if ((pred?.words || []).length === 0) {
    return pred.line;
  }
  const out: string[] = [pred.line.slice(0, pred.words[0].startPos)];
  pred.words.forEach((word: RskWordPosition, idx: number) => {
    if (idx > 0) {
      out.push(`<span class="not-matched">${pred.line.slice(pred?.words[idx-1].endPos, word.startPos)}</span>`);
    }
    out.push(`<span class="matched">${pred.line.slice(word.startPos, word.endPos)}</span>`);
  });
  out.push(`<span class="not-matched">${pred.line.slice(pred.words[pred.words.length-1].endPos, pred.line.length)}</span>`);

  return out.join("");
}

export function shortenStringToNearestWord(line: string, targetLength: number): string {
  if (line.length <= targetLength) {
    return line;
  }

  let out: string = "";
  line.split(" ").forEach((word: string) =>  {
    if (out.length > targetLength) {
      return;
    }
    out += ` ${word}`
  });
  return out
}
