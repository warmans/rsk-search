import { AfterViewInit, Component, Input, OnInit } from '@angular/core';
import { RsksearchDialog } from '../../../../lib/api-client/models';
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

  actorColorMap = {
    'ricky': '#fffdec',
    'steve': '#eeffec',
    'karl': '#ecffff',
  };

  constructor(private viewportScroller: ViewportScroller) {
  }

  ngOnInit(): void {
  }

  lineColor(d: RsksearchDialog): string {
    if (!d?.actor) {
      return 'transparent';
    }
    return this.actorColorMap[d.actor] || 'transparent';
  }

  ngAfterViewInit(): void {
    this.viewportScroller.scrollToAnchor(this.scrollToID);
  }
}
