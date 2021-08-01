import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TranscriptChangeComponent } from './page/transcript-change/transcript-change.component';
import { SharedModule } from '../shared/shared.module';
import { RouterModule } from '@angular/router';

@NgModule({
  declarations: [
    TranscriptChangeComponent
  ],
  imports: [
    CommonModule,
    RouterModule,
    SharedModule,
  ]
})
export class ContribModule {
}
