import { Component, Input, OnInit } from '@angular/core';
import { RskShortTranscript } from '../../../../lib/api-client/models';

@Component({
  selector: 'app-episode-summary',
  templateUrl: './episode-summary.component.html',
  styleUrls: ['./episode-summary.component.scss']
})
export class EpisodeSummaryComponent implements OnInit {

  @Input()
  episode: RskShortTranscript;

  constructor() {
  }

  ngOnInit(): void {
  }

}
