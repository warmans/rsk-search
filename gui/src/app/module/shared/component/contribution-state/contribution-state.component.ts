import { ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';
import { RskContributionState } from 'src/app/lib/api-client/models';
import { NgSwitch, NgIf, NgSwitchCase, NgSwitchDefault } from '@angular/common';

@Component({
  selector: 'app-contribution-state',
  templateUrl: './contribution-state.component.html',
  styleUrls: ['./contribution-state.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [NgSwitch, NgIf, NgSwitchCase, NgSwitchDefault],
})
export class ContributionStateComponent implements OnInit {
  @Input()
  state: RskContributionState;

  @Input()
  merged: boolean;

  states = RskContributionState;

  constructor() {}

  ngOnInit(): void {}
}
