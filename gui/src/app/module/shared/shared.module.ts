import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormErrorsComponent } from './component/form-errors/form-errors.component';
import { AlertComponent } from './component/alert/alert.component';
import { DropdownComponent } from './component/dropdown/dropdown.component';
import { FocusedDirective } from './directive/focused.directive';
import { LoadingOverlayComponent } from './component/loading-overlay/loading-overlay.component';
import { SafeUrlPipe } from './pipe/safe-url.pipe';
import { TranscriptComponent } from './component/transcript/transcript.component';
import { RouterModule } from '@angular/router';
import { AudioPlayerComponent } from './component/audio-player/audio-player.component';
import { SynopsesComponent } from './component/synopses/synopses.component';
import { WebhidDeviceComponent } from './component/webhid-device/webhid-device.component';
import { WebusbDeviceComponent } from './component/webusb-device/webusb-device.component';
import { ContributionStateComponent } from './component/contribution-state/contribution-state.component';

@NgModule({
  declarations: [
    FormErrorsComponent,
    AlertComponent,
    DropdownComponent,
    FocusedDirective,
    LoadingOverlayComponent,
    SafeUrlPipe,
    TranscriptComponent,
    AudioPlayerComponent,
    SynopsesComponent,
    WebhidDeviceComponent,
    WebusbDeviceComponent,
    ContributionStateComponent,
  ],
  imports: [
    CommonModule,
    RouterModule,
  ],
  providers: [],
  exports: [
    FormErrorsComponent,
    AlertComponent,
    DropdownComponent,
    LoadingOverlayComponent,
    SafeUrlPipe,
    TranscriptComponent,
    AudioPlayerComponent,
    SynopsesComponent,
    WebhidDeviceComponent,
    WebusbDeviceComponent,
    ContributionStateComponent,
  ]
})
export class SharedModule {
}
