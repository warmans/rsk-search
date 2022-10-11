import { ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';
import { RskSynopsis } from '../../../../lib/api-client/models';

@Component({
  selector: 'app-synopses',
  templateUrl: './synopses.component.html',
  styleUrls: ['./synopses.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SynopsesComponent implements OnInit {

  @Input()
  episodeID: string;

  @Input()
  synopses: RskSynopsis[];

  @Input()
  showTitle: boolean = true;

  constructor() {
  }

  ngOnInit(): void {
  }

}
