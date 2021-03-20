import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SubmitComponent } from './page/submit/submit.component';
import { EditorComponent } from './component/editor/editor.component';
import { SharedModule } from '../shared/shared.module';
import { RouterModule } from '@angular/router';


@NgModule({
  declarations: [SubmitComponent, EditorComponent],
  imports: [
    CommonModule,
    SharedModule,
    RouterModule
  ],
  exports: [SubmitComponent]
})
export class TscriptModule {
}
