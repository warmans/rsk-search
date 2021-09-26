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
    if (!value){
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
  scrollToID: string;

  @Input()
  enableAudioOffsets: boolean = false;

  @Input()
  enableLineLinks: boolean = false;

  @Output()
  emitAudioTimestamp: EventEmitter<number> = new EventEmitter();

  audioOffsetsAvailable: boolean = false;

  actorClassMap = {
    'ricky': 'ricky',
    'steve': 'steve',
    'karl': 'karl',
  };

  constructor(private viewportScroller: ViewportScroller) {
    viewportScroller.setOffset([0, 80]);
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
    this.viewportScroller.scrollToAnchor(this.scrollToID);
  }

  emitTimestamp(ts: number) {
    if (ts == 0) {
      return;
    }
    this.emitAudioTimestamp.next(ts);
  }
}
