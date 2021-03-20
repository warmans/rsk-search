import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SubmitComponent } from './page/submit/submit.component';
import { EditorComponent } from './component/editor/editor.component';
import { SharedModule } from '../shared/shared.module';


@NgModule({
  declarations: [SubmitComponent, EditorComponent],
  imports: [
    CommonModule,
    SharedModule
  ],
  exports: [SubmitComponent]
})
export class TscriptModule {
}
