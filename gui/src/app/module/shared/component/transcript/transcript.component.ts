import { AfterViewInit, Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { RskDialog } from '../../../../lib/api-client/models';
import { ViewportScroller } from '@angular/common';

@Component({
  selector: 'app-transcript',
  templateUrl: './transcript.component.html',
  styleUrls: ['./transcript.component.scss']
})
export class TranscriptComponent implements OnInit, AfterViewInit {


  @Input()
  set transcript(value: RskDialog[]) {
    if (!value) {
      return;
    }
    this._transcript = value;
    for (let i = 0; i < value.length; i++) {
      if (parseInt(value[i].offsetSec) > 0) {
        this.audioOffsetsAvailable = true;
        break;
      }
    }
  }

  get transcript(): RskDialog[] {
    return this._transcript;
  }

  private _transcript: RskDialog[];

  @Input()
  set scrollToID(value: string) {
    console.log('scroll to', value);
    this._scrollToID = value;
    this.scrollToAnchor();
  }

  get scrollToID(): string {
    return this._scrollToID;
  }

  private _scrollToID: string;

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
}
