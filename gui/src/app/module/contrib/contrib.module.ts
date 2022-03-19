import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TranscriptChangeComponent } from './page/transcript-change/transcript-change.component';
import { SharedModule } from '../shared/shared.module';
import { RouterModule } from '@angular/router';
import { EpisodeChunkSubmit } from './page/episode-chunk-submit/episode-chunk-submit.component';
import { RejectButtonComponent } from './component/reject-button/reject-button.component';
import { LeaderboardComponent } from './component/leaderboard/leaderboard.component';
import { EpisodeChunkContributions } from './page/episode-chunk-contributions/episode-chunk-contributions.component';
import { AuthorProfile } from './page/author-profile/author-profile.component';
import { ContributeComponent } from './page/contribute/contribute.component';
import { RandomComponent } from './page/random/random.component';
import { ReactiveFormsModule } from '@angular/forms';
import { RankingsComponent } from './component/rankings/rankings.component';

@NgModule({
  declarations: [
    TranscriptChangeComponent,
    EpisodeChunkSubmit,
    EpisodeChunkContributions,
    AuthorProfile,
    RandomComponent,
    ContributeComponent,
    RejectButtonComponent,
    LeaderboardComponent,
    RankingsComponent
  ],
  imports: [
    CommonModule,
    RouterModule,
    SharedModule,
    ReactiveFormsModule,
  ]
})
export class ContribModule {
}
