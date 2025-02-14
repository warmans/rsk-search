import { RskPrediction } from './api-client/models';

export function trimChars(str: string, ch: string): string {
  let start = 0,
    end = str.length;

  while (start < end && str[start] === ch) {
    ++start;
  }
  while (end > start && str[end - 1] === ch) {
    --end;
  }
  return start > 0 || end < str.length ? str.substring(start, end) : str;
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

export function formatSecondsAsTimestamp(seconds: number | string, secondsAreMilliseconds?: boolean): string {
  if (!seconds) {
    return '-';
  }

  let secondsNum: number = typeof seconds === 'string' ? parseInt(seconds) : seconds;
  if (secondsAreMilliseconds) {
    secondsNum = secondsNum / 1000;
  }
  const minsNum: number = secondsNum / 60;
  const mins: string = String(Math.floor(minsNum).toFixed(0)).padStart(2, '0');
  const secs: string = String(((minsNum % 1) * 60).toFixed(0)).padStart(2, '0');
  return `${mins}:${secs}`;
}

export function episodeIdVariations(id: string): string[] {
  const [matches] = Array.from(id.matchAll(new RegExp('([a-z]+)-S([0-9]+)E([0-9]+)$', 'g')));
  if (!matches) {
    return [id];
  }
  const [, pub, series, ep] = matches;
  return [id, `${pub}-S${parseInt(series)}E${ep}`, `${pub}-S0${parseInt(series)}E${ep}`];
}
