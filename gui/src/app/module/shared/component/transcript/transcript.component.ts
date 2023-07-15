import {AfterViewInit, ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {RskDialog, RskSynopsis, RskTranscript, RskTrivia} from '../../../../lib/api-client/models';
import {ViewportScroller} from '@angular/common';
import {parseTranscript, Tscript} from '../../lib/tscript';
import {ClipboardService} from 'src/app/module/core/service/clipboard/clipboard.service';

interface DialogGroup {
  startPos: number;
  endPos: number;
  tscript: Tscript;
}

export interface Section {
  epid: string;
  startPos: number;
  endPos: number;
}

@Component({
  selector: 'app-transcript',
  templateUrl: './transcript.component.html',
  styleUrls: ['./transcript.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class TranscriptComponent implements OnInit, AfterViewInit {

  @Input()
  epid: string;

  @Input()
  set transcript(value: Tscript | RskTranscript) {
    if (!value) {
      return;
    }
    this._transcript = value;
    this.preProcessTranscript(value);
  }

  get transcript(): Tscript | RskTranscript {
    return this._transcript;
  }

  private _transcript: Tscript | RskTranscript;

  @Input()
  set rawTranscript(value: string) {
    this._rawTranscript = value;
    this.transcript = parseTranscript(value);
  }

  get rawTranscript(): string {
    return this._rawTranscript;
  }

  private _rawTranscript: string;

  groupedDialog: DialogGroup[];

  lineInSynopsisMap: { [index: number]: boolean } = {};
  synopsisPos: { [index: number]: RskSynopsis } = {};

  @Input()
  set scrollToID(value: string | null) {
    if (value === null) {
      return;
    }
    this._scrollToID = value;
    this.scrollToAnchor();
  }

  get scrollToID(): string {
    return this._scrollToID;
  }

  private _scrollToID: string;

  @Input()
  set scrollToSeconds(value: number | null) {
    if (value === null) {
      return;
    }
    this._scrollToSeconds = value;
    this._scrollToSecondOffset(value);
  }

  get scrollToSeconds(): number {
    return this._scrollToSeconds;
  }

  private _scrollToSeconds: number;

  scrollToPosStart: number;
  scrollToPosEnd: number;

  @Input()
  searchResultMode: boolean = false;

  @Input()
  enableLineLinking: boolean = false;

  @Input()
  enableLineCopy: boolean = false;

  @Input()
  enableAudioLinks: boolean = true;

  @Input()
  enableShareLinks: boolean = true;

  @Input()
  startLine: number;

  @Input()
  endLine: number;

  @Output()
  emitAudioTimestamp: EventEmitter<number> = new EventEmitter();

  @Output()
  emitShare: EventEmitter<Section> = new EventEmitter();

  @Output()
  emitSelection: EventEmitter<Section> = new EventEmitter();

  audioOffsetsAvailable: boolean = false;

  actorClassMap = {
    'ricky': 'ricky',
    'steve': 'steve',
    'karl': 'karl',
    'claire': 'claire',
    'camfield': 'camfield',
  };

  constructor(private viewportScroller: ViewportScroller, private clipboard: ClipboardService) {
    viewportScroller.setOffset([0, window.innerHeight / 2]);
  }

  ngOnInit(): void {
  }

  actorClass(d: RskDialog): string {
    if (!d?.actor) {
      return '';
    }
    return this.actorClassMap[d.actor.toLowerCase().trim()] || '';
  }

  ngAfterViewInit(): void {
    this.scrollToAnchor();
  }

  scrollToAnchor() {
    if (!this._scrollToID) {
      this.scrollToPosStart = undefined;
      this.scrollToPosEnd = undefined;
      return;
    }
    let parts = this._scrollToID.split('-');
    if (parts.length === 2) {
      this.viewportScroller.scrollToAnchor(this._scrollToID);
      this.scrollToPosStart = parseInt(parts[1]);
      this.scrollToPosEnd = this.scrollToPosStart;
    } else if (parts.length === 3) {
      this.viewportScroller.scrollToAnchor(`${parts[0]}-${parts[1]}`);
      this.scrollToPosStart = parseInt(parts[1]);
      this.scrollToPosEnd = parseInt(parts[2]);
    }
  }

  private _scrollToSecondOffset(seconds: number) {
    for (let i = 0; i < this.transcript.transcript.length; i++) {
      if (parseInt(this.transcript.transcript[i].offsetSec) >= seconds) {
        this.scrollToID = `pos-${this.transcript.transcript[i].pos}`;
        return;
      }
    }
  }

  emitTimestamp(ts: string) {
    const tsInt = parseInt(ts);
    if (!tsInt) {
      return;
    }
    this.emitAudioTimestamp.next(tsInt);
  }

  selectPosition(pos: number, ev: any): boolean {
    if (ev.shiftKey && this.scrollToPosStart) {
      const start = this.scrollToPosStart > pos ? pos : this.scrollToPosStart;
      const end = this.scrollToPosStart > pos ? this.scrollToPosStart : pos;
      this.emitSelection.next({startPos: start, endPos: end, epid: this.epid});
      return false;
    }
    this.emitSelection.next({startPos: pos, endPos: pos, epid: this.epid});
    return true;
  }

  selectRange(startLine: number, endLine: number): boolean {
    this.emitSelection.next({startPos: startLine, endPos: endLine, epid: this.epid});
    return true;
  }

  preProcessTranscript(episode: Tscript | RskTranscript) {

    if (!episode) {
      return;
    }

    this.synopsisPos = {};
    (this.transcript?.synopses || []).forEach((s) => {
      this.synopsisPos[s.startPos] = s;
    });

    this.lineInSynopsisMap = {};
    this.audioOffsetsAvailable = false;
    this.groupedDialog = [];
    let currentGroup: DialogGroup = {
      startPos: 1,
      endPos: undefined,
      tscript: {synopses: [], trivia: [], transcript: []}
    };
    for (let i = (this.startLine || 0); i < (this.endLine && this.endLine < episode?.transcript.length ? this.endLine : episode?.transcript.length); i++) {

      if (parseInt(episode.transcript[i].offsetSec) > 0) {
        this.audioOffsetsAvailable = true;
      }

      this.lineInSynopsisMap[episode.transcript[i].pos] = !!(episode?.synopses || []).find((s: RskSynopsis) => episode.transcript[i].pos >= s.startPos && i < s.endPos);

      currentGroup.tscript.transcript.push(episode.transcript[i]);

      // there may be multiple trivias which intersect this line. So find them all and then,
      // append them as required.
      const foundIntersectingTrivia: RskTrivia[] = (episode?.trivia || []).filter((s: RskSynopsis) => episode.transcript[i].pos === s.startPos - 1 || episode.transcript[i].pos === s.endPos);

      if ((foundIntersectingTrivia || []).length > 0) {
        foundIntersectingTrivia.forEach((trivia: RskTrivia) => {
          if (episode.transcript[i].pos === trivia.startPos - 1) {

            // flush current group
            currentGroup.endPos = episode.transcript[i].pos;
            this.groupedDialog.push(currentGroup);

            // start a new group
            currentGroup = {
              startPos: episode.transcript[i].pos,
              endPos: undefined,
              tscript: {synopses: [], trivia: [], transcript: []}
            };
            if ((currentGroup?.tscript?.trivia || []).length === 0) {
              currentGroup.tscript.trivia = [trivia];
            } else {
              currentGroup.tscript.trivia.push(trivia);
            }
          }
          if (episode.transcript[i].pos === trivia.endPos) {
            // flush current group
            currentGroup.endPos = episode.transcript[i].pos;
            this.groupedDialog.push(currentGroup);

            // start a new group
            currentGroup = {
              startPos: episode.transcript[i].pos,
              endPos: undefined,
              tscript: {synopses: [], trivia: [], transcript: []}
            };
          }
        })
      }
    }
    if (currentGroup.endPos === undefined && (episode.transcript || []).length > 0) {
      currentGroup.endPos = episode.transcript[episode.transcript.length - 1].pos;
      this.groupedDialog.push(currentGroup);
    }
  }

  emitShareOpts(startPos: number, endPos: number) {
    this.emitShare.next({epid: this.epid, startPos: startPos, endPos: endPos});
  }

  copyLineToClipboard(content: string, timestamp?: number) {
    this.clipboard.copyTextToClipboard(content);
  }
}
