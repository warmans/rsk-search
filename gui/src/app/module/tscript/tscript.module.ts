import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SubmitComponent } from './page/submit/submit.component';
import { SharedModule } from '../shared/shared.module';
import { RouterModule } from '@angular/router';
import { RandomComponent } from './page/random/random.component';
import { AuthorContributionsComponent } from './page/author-contributions/author-contributions.component';
import { ApproveComponent } from './page/approve/approve.component';
import { ContributeComponent } from './page/contribute/contribute.component';
import { LeaderboardComponent } from './component/leaderboard/leaderboard.component';
import { RewardModule } from '../reward/reward.module';
import { SubmitV2Component } from './page/submit-v2/submit.component';
import { RejectButtonComponent } from './component/reject-button/reject-button.component';
import { ReactiveFormsModule } from '@angular/forms';

@NgModule({
  declarations: [
    SubmitComponent,
    SubmitV2Component,
    RandomComponent,
    AuthorContributionsComponent,
    ApproveComponent,
    ContributeComponent,
    LeaderboardComponent,
    RejectButtonComponent,
  ],
  imports: [
    CommonModule,
    SharedModule,
    RouterModule,
    RewardModule,
    ReactiveFormsModule,
  ],
  exports: []
})
export class TscriptModule {
}
