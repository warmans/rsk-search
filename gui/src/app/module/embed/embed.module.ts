import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TranscriptSectionComponent } from './page/transcript-section/transcript-section.component';
import { SharedModule } from '../shared/shared.module';

@NgModule({
  declarations: [
    TranscriptSectionComponent
  ],
  imports: [
    CommonModule,
    SharedModule,
  ]
})
export class EmbedModule { }
