import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TranscriptSectionComponent } from './page/transcript-section/transcript-section.component';
import { SharedModule } from '../shared/shared.module';
import { PreviewModalComponent } from './component/preview-modal/preview-modal.component';
import { ReactiveFormsModule } from '@angular/forms';

@NgModule({
  declarations: [
    TranscriptSectionComponent,
    PreviewModalComponent
  ],
  exports: [
    PreviewModalComponent
  ],
  imports: [
    CommonModule,
    SharedModule,
    ReactiveFormsModule,
  ]
})
export class EmbedModule {
}
