import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormErrorsComponent } from './component/form-errors/form-errors.component';
import { AlertComponent } from './component/alert/alert.component';

@NgModule({
  declarations: [FormErrorsComponent, AlertComponent],
  imports: [
    CommonModule,
  ],
  providers: [],
  exports: [
    FormErrorsComponent,
    AlertComponent,
  ]
})
export class SharedModule {
}
