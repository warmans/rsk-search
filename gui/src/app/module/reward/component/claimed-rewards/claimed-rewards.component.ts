import { Component, Input, OnInit } from '@angular/core';
import { RskClaimedReward } from '../../../../lib/api-client/models';
import { environment } from '../../../../../environments/environment.prod';

@Component({
  selector: 'app-claimed-rewards',
  templateUrl: './claimed-rewards.component.html',
  styleUrls: ['./claimed-rewards.component.scss']
})
export class ClaimedRewardsComponent implements OnInit {

  @Input()
  set rewards(value: RskClaimedReward[]) {
    if (!value) {
      return;
    }
    this._rewards = value;
    value.forEach((row) => {
      this.totalValue += row.claimValue;
    });
  }

  get rewards(): RskClaimedReward[] {
    return this._rewards;
  }

  private _rewards: RskClaimedReward[] = [];

  totalValue: number = 0;

  awardThreshold = environment.contributionAwardThreshold;

  constructor() {
  }

  ngOnInit(): void {
  }

}
