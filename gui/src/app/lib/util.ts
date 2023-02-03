import { RskPrediction } from './api-client/models';

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
  if (pred.fragment === '') {
    return pred.line;
  }
  let out: string = '';
  out = pred.fragment.replace(/\{\{/g, '<span class="matched">');
  out = out.replace(/}}/g, '</span>');
  return out;
}

export function shortenStringToNearestWord(line: string, targetLength: number): string {
  if (line.length <= targetLength) {
    return line;
  }

  let out: string = '';
  line.split(' ').forEach((word: string) => {
    if (out.length > targetLength) {
      return;
    }
    out += ` ${word}`;
  });
  return out;
}
