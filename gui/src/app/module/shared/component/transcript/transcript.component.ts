import { AfterViewInit, ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { RskDialog, RskSynopsis } from '../../../../lib/api-client/models';
import { ViewportScroller } from '@angular/common';
import { parseTranscript, Tscript } from '../../lib/tscript';

interface DialogGroup {
  startPos: number;
  endPos: number;

  tscript: Tscript;
}

@Component({
  selector: 'app-transcript',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './transcript.component.html',
  styleUrls: ['./transcript.component.scss']
})
export class TranscriptComponent implements OnInit, AfterViewInit {

  @Input()
  set transcript(value: Tscript) {
    if (!value) {
      return;
    }
    this._transcript = value;
    this.preProcessTranscript(value);
  }

  get transcript(): Tscript {
    return this._transcript;
  }

  private _transcript: Tscript;

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
  synopsisTitlePos: { [index: number]: string } = {};

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

  @Output()
  emitAudioTimestamp: EventEmitter<number> = new EventEmitter();

  audioOffsetsAvailable: boolean = false;

  actorClassMap = {
    'ricky': 'ricky',
    'steve': 'steve',
    'karl': 'karl',
    'claire': 'claire',
    'camfield': 'camfield',
  };

  constructor(private viewportScroller: ViewportScroller) {
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

  clearFocus() {
    this.scrollToPosStart = undefined;
    this.scrollToPosEnd = undefined;
  }

  emitTimestamp(ts: number) {
    if (ts == 0) {
      return;
    }
    this.emitAudioTimestamp.next(ts);
  }

  preProcessTranscript(value: Tscript) {

    if (!value) {
      return;
    }

    this.synopsisTitlePos = {};
    (this.transcript?.synopses || []).forEach((s) => {
      this.synopsisTitlePos[s.startPos] = s.description;
    });

    this.lineInSynopsisMap = {};
    this.audioOffsetsAvailable = false;
    this.groupedDialog = [];
    let currentGroup: DialogGroup = {
      startPos: 1,
      endPos: undefined,
      tscript: { synopses: [], trivia: [], transcript: [] }
    };
    for (let i = 0; i < value?.transcript.length; i++) {

      if (parseInt(value.transcript[i].offsetSec) > 0) {
        this.audioOffsetsAvailable = true;
      }

      this.lineInSynopsisMap[value.transcript[i].pos] = !!(value?.synopses || []).find((s: RskSynopsis) => value.transcript[i].pos >= s.startPos && i < s.endPos);

      currentGroup.tscript.transcript.push(value.transcript[i]);

      const foundIntersectingTrivia = (value?.trivia || []).find((s: RskSynopsis) => value.transcript[i].pos === s.startPos - 1 || value.transcript[i].pos === s.endPos);
      if (foundIntersectingTrivia) {
        if (value.transcript[i].pos === foundIntersectingTrivia.startPos - 1) {
          currentGroup.endPos = value.transcript[i].pos;
          this.groupedDialog.push(currentGroup);
          currentGroup = {
            startPos: value.transcript[i].pos,
            endPos: undefined,
            tscript: { synopses: [], trivia: [foundIntersectingTrivia], transcript: [] }
          };
        }
        if (value.transcript[i].pos === foundIntersectingTrivia.endPos) {
          currentGroup.endPos = value.transcript[i].pos;
          this.groupedDialog.push(currentGroup);
          currentGroup = {
            startPos: value.transcript[i].pos,
            endPos: undefined,
            tscript: { synopses: [], trivia: [], transcript: [] }
          };
        }
      }
    }
    if (currentGroup.endPos === undefined) {
      currentGroup.endPos = value.transcript.length;
      this.groupedDialog.push(currentGroup);
    }
  }
}
