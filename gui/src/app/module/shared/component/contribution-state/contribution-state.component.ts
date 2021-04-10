import { Component, Input, OnInit } from '@angular/core';
import { RsksearchContributionState } from '../../../../lib/api-client/models';

@Component({
  selector: 'app-contribution-state',
  templateUrl: './contribution-state.component.html',
  styleUrls: ['./contribution-state.component.scss']
})
export class ContributionStateComponent implements OnInit {

  @Input()
  state: RsksearchContributionState;

  states = RsksearchContributionState;

  constructor() {
  }

  ngOnInit(): void {
  }

}
