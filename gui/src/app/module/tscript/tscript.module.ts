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

@NgModule({
  declarations: [
    SubmitComponent,
    RandomComponent,
    AuthorContributionsComponent,
    ApproveComponent,
    ContributeComponent,
    LeaderboardComponent,
  ],
  imports: [
    CommonModule,
    SharedModule,
    RouterModule,
    RewardModule,
  ],
  exports: [SubmitComponent]
})
export class TscriptModule {
}
