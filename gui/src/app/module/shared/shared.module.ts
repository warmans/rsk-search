import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormErrorsComponent } from './component/form-errors/form-errors.component';
import { AlertComponent } from './component/alert/alert.component';
import { FocusedDirective } from './directive/focused.directive';
import { LoadingOverlayComponent } from './component/loading-overlay/loading-overlay.component';
import { SafeUrlPipe } from './pipe/safe-url.pipe';
import { TranscriptComponent } from './component/transcript/transcript.component';
import { RouterModule } from '@angular/router';
import { SynopsesComponent } from './component/synopses/synopses.component';
import { ContributionStateComponent } from './component/contribution-state/contribution-state.component';
import { RedditLoginComponent } from './component/reddit-login/reddit-login.component';
import { EditorComponent } from './component/editor/editor.component';
import { EditorConfigComponent } from './component/editor-config/editor-config.component';
import { ReactiveFormsModule } from '@angular/forms';
import { EditorHelpComponent } from './component/editor-help/editor-help.component';
import { TranscriberComponent } from './component/transcriber/transcriber.component';
import { FormatSecondsPipe } from './pipe/format-seconds.pipe';
import { TriviaComponent } from './component/trivia/trivia.component';
import { HtmlDiffComponent } from './component/html-diff/html-diff.component';
import { ClaimedRewardsComponent } from './component/claimed-rewards/claimed-rewards.component';
import { MarkdownComponent } from './component/markdown/markdown.component';
import { AudioPlayerV2Component } from './component/audio-player-v2/audio-player-v2.component';
import { AudioPlayerFixedComponent } from './component/audio-player-fixed/audio-player-fixed.component';
import { NgbPopoverModule } from '@ng-bootstrap/ng-bootstrap';
import { FindReplaceComponent } from './component/find-replace/find-replace.component';

@NgModule({
  declarations: [
    FormErrorsComponent,
    AlertComponent,
    FocusedDirective,
    LoadingOverlayComponent,
    SafeUrlPipe,
    FormatSecondsPipe,
    TranscriptComponent,
    SynopsesComponent,
    TriviaComponent,
    ContributionStateComponent,
    RedditLoginComponent,
    EditorComponent,
    EditorConfigComponent,
    EditorHelpComponent,
    TranscriberComponent,
    HtmlDiffComponent,
    ClaimedRewardsComponent,
    MarkdownComponent,
    AudioPlayerV2Component,
    AudioPlayerFixedComponent,
    FindReplaceComponent,
  ],
  imports: [
    CommonModule,
    RouterModule,
    ReactiveFormsModule,
    NgbPopoverModule,
  ],
  providers: [],
  exports: [
    FormErrorsComponent,
    AlertComponent,
    LoadingOverlayComponent,
    SafeUrlPipe,
    FormatSecondsPipe,
    TranscriptComponent,
    SynopsesComponent,
    TriviaComponent,
    ContributionStateComponent,
    RedditLoginComponent,
    EditorComponent,
    EditorConfigComponent,
    EditorHelpComponent,
    TranscriberComponent,
    ClaimedRewardsComponent,
    MarkdownComponent,
    AudioPlayerV2Component,
    AudioPlayerFixedComponent,
  ]
})
export class SharedModule {
}
