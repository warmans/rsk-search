import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormErrorsComponent} from './component/form-errors/form-errors.component';
import {AlertComponent} from './component/alert/alert.component';
import {FocusedDirective} from './directive/focused.directive';
import {LoadingOverlayComponent} from './component/loading-overlay/loading-overlay.component';
import {SafeUrlPipe} from './pipe/safe-url.pipe';
import {TranscriptComponent} from './component/transcript/transcript.component';
import {RouterModule} from '@angular/router';
import {SynopsesComponent} from './component/synopses/synopses.component';
import {ContributionStateComponent} from './component/contribution-state/contribution-state.component';
import {OauthLoginComponent} from './component/oauth-login/oauth-login.component';
import {EditorInputComponent} from './component/editor-input/editor-input.component';
import {EditorConfigComponent} from './component/editor-config/editor-config.component';
import {ReactiveFormsModule} from '@angular/forms';
import {EditorHelpComponent} from './component/editor-help/editor-help.component';
import {EditorComponent} from './component/editor/editor.component';
import {FormatSecondsPipe} from './pipe/format-seconds.pipe';
import {TriviaComponent} from './component/trivia/trivia.component';
import {HtmlDiffComponent} from './component/html-diff/html-diff.component';
import {ClaimedRewardsComponent} from './component/claimed-rewards/claimed-rewards.component';
import {MarkdownComponent} from './component/markdown/markdown.component';
import {AudioPlayerV2Component} from './component/audio-player-v2/audio-player-v2.component';
import {AudioPlayerFixedComponent} from './component/audio-player-fixed/audio-player-fixed.component';
import {NgbPopoverModule} from '@ng-bootstrap/ng-bootstrap';
import {FindReplaceComponent} from './component/find-replace/find-replace.component';
import {MetadataEditorComponent} from './component/metadata-editor/metadata-editor.component';
import {UserMenuComponent} from './component/user-menu/user-menu.component';
import {TranscriptCopyDialogComponent} from './component/transcript-copy-dialog/transcript-copy-dialog.component';
import {RandomQuoteComponent} from "./component/random-quote/random-quote.component";

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
    OauthLoginComponent,
    EditorInputComponent,
    EditorConfigComponent,
    EditorHelpComponent,
    EditorComponent,
    HtmlDiffComponent,
    ClaimedRewardsComponent,
    MarkdownComponent,
    AudioPlayerV2Component,
    AudioPlayerFixedComponent,
    FindReplaceComponent,
    MetadataEditorComponent,
    UserMenuComponent,
    TranscriptCopyDialogComponent,
    RandomQuoteComponent,
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
        OauthLoginComponent,
        EditorInputComponent,
        EditorConfigComponent,
        EditorHelpComponent,
        EditorComponent,
        ClaimedRewardsComponent,
        MarkdownComponent,
        AudioPlayerV2Component,
        AudioPlayerFixedComponent,
        HtmlDiffComponent,
        MetadataEditorComponent,
        UserMenuComponent,
        TranscriptCopyDialogComponent,
        RandomQuoteComponent,
    ]
})
export class SharedModule {
}
