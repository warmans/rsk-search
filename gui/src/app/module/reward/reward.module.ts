import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RedeemComponent } from './page/redeem/redeem.component';
import { PendingRewardsComponent } from './component/pending-rewards/pending-rewards.component';
import { RouterModule } from '@angular/router';
import { ReactiveFormsModule } from '@angular/forms';
import { SharedModule } from '../shared/shared.module';
import { ClaimedRewardsComponent } from './component/claimed-rewards/claimed-rewards.component';

@NgModule({
  declarations: [RedeemComponent, PendingRewardsComponent, ClaimedRewardsComponent],
  imports: [
    CommonModule,
    RouterModule,
    ReactiveFormsModule,
    SharedModule,
  ],
  exports: [RedeemComponent, PendingRewardsComponent, ClaimedRewardsComponent]
})
export class RewardModule {
}
