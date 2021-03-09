import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormErrorsComponent } from './component/form-errors/form-errors.component';
import { AlertComponent } from './component/alert/alert.component';
import { DropdownComponent } from './component/dropdown/dropdown.component';
import { FocusedDirective } from './directive/focused.directive';

@NgModule({
  declarations: [FormErrorsComponent, AlertComponent, DropdownComponent, FocusedDirective],
  imports: [
    CommonModule,
  ],
  providers: [],
  exports: [
    FormErrorsComponent,
    AlertComponent,
    DropdownComponent,
  ]
})
export class SharedModule {
}
