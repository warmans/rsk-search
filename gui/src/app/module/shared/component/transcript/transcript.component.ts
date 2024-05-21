import {AfterViewInit, ChangeDetectionStrategy, Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {
  DialogType,
  RskDialog,
  RskMediaType,
  RskSynopsis,
  RskTranscript,
  RskTrivia
} from '../../../../lib/api-client/models';
import {ViewportScroller} from '@angular/common';
import {parseTranscript, Tscript} from '../../lib/tscript';
import {ClipboardService} from 'src/app/module/core/service/clipboard/clipboard.service';
import {BehaviorSubject, Subject} from "rxjs";
import {takeUntil} from "rxjs/operators";

interface DialogGroup {
  startPos: number;
  endPos: number;
  tscript: Tscript;
}

export interface Section {
  epid?: string;
  startPos: number;
  startTimestampMs?: number;
  endPos?: number;
  endTimestampMs?: number;
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
    this._transcript.next(value);
  }

  get transcript(): Tscript | RskTranscript {
    return this._transcript.value;
  }

  private _transcript: BehaviorSubject<Tscript | RskTranscript> = new BehaviorSubject<Tscript | RskTranscript>(null);

  @Input()
  set rawTranscript(value: string) {
    this._rawTranscript = value;
    this.transcript = parseTranscript(value);
  }

  get rawTranscript(): string {
    return this._rawTranscript;
  }

  private _rawTranscript: string;

  @Input()
  mediaType: RskMediaType = RskMediaType.AUDIO;

  groupedDialog: DialogGroup[];

  dialogTypes = DialogType;

  lineInSynopsisMap: { [index: number]: boolean } = {};
  synopsisPos: { [index: number]: RskSynopsis } = {};

  idScrollerSubject: Subject<string | null> = new BehaviorSubject(null)
  scrollAnchor: string;

  @Input()
  set scrollToID(value: string | null) {
    if (value === null) {
      return;
    }
    this._scrollToID = value;

    if (!this._scrollToID) {
      this.scrollToPosStart = undefined;
      this.scrollToPosEnd = undefined;
      return;
    }
    let parts = this._scrollToID.split('-');
    if (parts.length === 2) {
      this.scrollAnchor = this._scrollToID;
      this.scrollToPosStart = parseInt(parts[1]);
      this.scrollToPosEnd = this.scrollToPosStart;
    } else if (parts.length === 3) {
      this.scrollAnchor = `${parts[0]}-${parts[1]}`;
      this.scrollToPosStart = parseInt(parts[1]);
      this.scrollToPosEnd = parseInt(parts[2]);
    }
    this.idScrollerSubject.next(value);
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
  startLine: number;

  @Input()
  endLine: number;

  @Output()
  emitAudioTimestamp: EventEmitter<number> = new EventEmitter();

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

  destroy$: Subject<void> = new Subject<void>();

  constructor(private viewportScroller: ViewportScroller, private clipboard: ClipboardService) {
    viewportScroller.setOffset([0, window.innerHeight / 2]);
  }

  ngOnInit(): void {
    this._transcript.pipe(takeUntil(this.destroy$)).subscribe((transcript) => {
      if (transcript === null) {
        return
      }
      this.preProcessTranscript(transcript);
    });
  }

  actorClass(d: RskDialog): string {
    if (!d?.actor) {
      return '';
    }
    return this.actorClassMap[d.actor.toLowerCase().trim()] || '';
  }

  ngAfterViewInit(): void {
    this.idScrollerSubject.pipe(takeUntil(this.destroy$)).subscribe((id: string) => {
      if (id == null) {
        return;
      }
      this.scrollToAnchor();
    })
  }

  scrollToAnchor() {
    this.viewportScroller.scrollToAnchor(this.scrollAnchor);
  }

  private _scrollToSecondOffset(seconds: number) {
    for (let i = 0; i < this.transcript.transcript.length; i++) {
      if (this.transcript.transcript[i].offsetMs / 1000 >= seconds) {
        this.scrollToID = `pos-${this.transcript.transcript[i].pos}`;
        return;
      }
    }
  }

  emitTimestamp(ts: number) {
    this.emitAudioTimestamp.next(ts);
  }

  selectPosition(line: RskDialog, ev: any): boolean {
    this.emitSelection.next({
      startPos: line.pos,
      startTimestampMs: line.offsetMs,
      endPos: line.pos,
      endTimestampMs: (line.offsetMs + line.durationMs),
      epid: this.epid,
    });
    return true;
  }

  addToSelection(line: RskDialog): boolean {
    if (this.scrollToPosStart) {
      const start = this.scrollToPosStart > line.pos ? line.pos : this.scrollToPosStart;
      const end = this.scrollToPosStart > line.pos ? this.scrollToPosStart : line.pos;
      this.emitSelection.next({
        startPos: start,
        startTimestampMs: line.offsetMs,
        endPos: end,
        endTimestampMs: line.offsetMs + line.durationMs,
        epid: this.epid,
      });
    }
    return false;
  }

  selectRange(startPos: number, endPos: number): boolean {
    this.emitSelection.next({
      startPos: startPos,
      startTimestampMs: this.transcript.transcript[startPos - 1].durationMs,
      endPos: endPos,
      endTimestampMs: this.transcript.transcript[endPos - 1]?.offsetMs + this.transcript.transcript[endPos - 1]?.durationMs,
      epid: this.epid,
    });
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

    // add a fake trivia for song information
    if (this.searchResultMode === false) {
      episode.transcript.forEach((dialog) => {
        if (dialog.type == DialogType.SONG && dialog.metadata && dialog.metadata["song_album_art"]) {
          if (!episode.trivia) {
            episode.trivia = [];
          }
          episode.trivia.push({
            description: `<img src="${dialog.metadata["song_album_art"]}" alt="${dialog.metadata["song_album"]}" width="300px"/>`,
            startPos: dialog.pos,
            endPos: Math.min(dialog.pos + 5, episode.transcript.length),
          })
        }
      });
    }

    this.lineInSynopsisMap = {};
    this.audioOffsetsAvailable = false;
    this.groupedDialog = [];
    let currentGroup: DialogGroup = {
      startPos: 1,
      endPos: undefined,
      tscript: {synopses: [], trivia: [], transcript: []}
    };
    for (let i: number = (this.startLine || 0); i < (this.endLine && this.endLine < episode?.transcript.length ? this.endLine : episode?.transcript.length); i++) {

      if (episode.transcript[i].offsetMs > 0) {
        this.audioOffsetsAvailable = true;
      }

      this.lineInSynopsisMap[episode.transcript[i].pos] = !!(episode?.synopses || []).find((s: RskSynopsis): boolean => episode.transcript[i].pos >= s.startPos && i < s.endPos);

      currentGroup.tscript.transcript.push(episode.transcript[i]);

      // there may be multiple trivias which intersect this line. So find them all and then,
      // append them as required.
      const foundIntersectingTrivia: RskTrivia[] = (episode?.trivia || []).filter((s: RskTrivia) => episode.transcript[i].pos === s.startPos - 1 || episode.transcript[i].pos === s.endPos);

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

  copyLineToClipboard(content: string, timestamp?: number) {
    this.clipboard.copyTextToClipboard(content);
  }

  protected readonly RskMediaType = RskMediaType;
}
