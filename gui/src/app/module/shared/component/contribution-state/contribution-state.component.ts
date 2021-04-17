import { Component, Input, OnInit } from '@angular/core';
import { RskContributionState } from '../../../../lib/api-client/models';

@Component({
  selector: 'app-contribution-state',
  templateUrl: './contribution-state.component.html',
  styleUrls: ['./contribution-state.component.scss']
})
export class ContributionStateComponent implements OnInit {

  @Input()
  state: RskContributionState;

  states = RskContributionState;

  constructor() {
  }

  ngOnInit(): void {
  }

}
