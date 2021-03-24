import { AfterViewInit, Component, Input, OnInit } from '@angular/core';
import { RsksearchDialog, RsksearchSynopsis } from '../../../../lib/api-client/models';
import { ViewportScroller } from '@angular/common';

@Component({
  selector: 'app-transcript',
  templateUrl: './transcript.component.html',
  styleUrls: ['./transcript.component.scss']
})
export class TranscriptComponent implements OnInit, AfterViewInit {

  @Input()
  transcript: RsksearchDialog[];

  @Input()
  scrollToID: string;

  @Input()
  enableLineLinks: boolean = false;

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

  actorClass(d: RsksearchDialog): string {
    if (!d?.actor) {
      return '';
    }
    return this.actorClassMap[d.actor] || '';
  }


  ngAfterViewInit(): void {
    this.viewportScroller.scrollToAnchor(this.scrollToID);
  }
}
