import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SubmitComponent } from './page/submit/submit.component';
import { EditorComponent } from './component/editor/editor.component';
import { SharedModule } from '../shared/shared.module';
import { RouterModule } from '@angular/router';
import { RandomComponent } from './page/random/random.component';
import { AuthorContributionsComponent } from './page/author-contributions/author-contributions.component';
import { RedditLoginComponent } from './component/reddit-login/reddit-login.component';
import { ReactiveFormsModule } from '@angular/forms';
import { ApproveComponent } from './page/approve/approve.component';


@NgModule({
  declarations: [SubmitComponent, EditorComponent, RandomComponent, AuthorContributionsComponent, RedditLoginComponent, ApproveComponent],
  imports: [
    CommonModule,
    SharedModule,
    RouterModule,
    ReactiveFormsModule,
  ],
  exports: [SubmitComponent]
})
export class TscriptModule {
}
