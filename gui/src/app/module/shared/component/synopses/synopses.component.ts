import { Component, Input, OnInit } from '@angular/core';
import { RsksearchSynopsis } from '../../../../lib/api-client/models';

@Component({
  selector: 'app-synopses',
  templateUrl: './synopses.component.html',
  styleUrls: ['./synopses.component.scss']
})
export class SynopsesComponent implements OnInit {

  @Input()
  synopses: RsksearchSynopsis[];

  constructor() {
  }

  ngOnInit(): void {
  }

}
