import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SubmitComponent } from './page/submit/submit.component';
import { EditorComponent } from './component/editor/editor.component';
import { SharedModule } from '../shared/shared.module';
import { RouterModule } from '@angular/router';
import { RandomComponent } from './page/random/random.component';
import { AuthorContributionsComponent } from './page/author-contributions/author-contributions.component';
import { RedditLoginComponent } from './component/reddit-login/reddit-login.component';
import { ContenteditableModule } from '@ng-stack/contenteditable';
import { ReactiveFormsModule } from '@angular/forms';


@NgModule({
  declarations: [SubmitComponent, EditorComponent, RandomComponent, AuthorContributionsComponent, RedditLoginComponent],
  imports: [
    CommonModule,
    SharedModule,
    RouterModule,
    ContenteditableModule,
    ReactiveFormsModule,
  ],
  exports: [SubmitComponent]
})
export class TscriptModule {
}
