import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormErrorsComponent } from './component/form-errors/form-errors.component';
import { AlertComponent } from './component/alert/alert.component';
import { DropdownComponent } from './component/dropdown/dropdown.component';
import { FocusedDirective } from './directive/focused.directive';
import { LoadingOverlayComponent } from './component/loading-overlay/loading-overlay.component';
import { SafeUrlPipe } from './pipe/safe-url.pipe';

@NgModule({
  declarations: [
    FormErrorsComponent,
    AlertComponent,
    DropdownComponent,
    FocusedDirective,
    LoadingOverlayComponent,
    SafeUrlPipe
  ],
  imports: [
    CommonModule,
  ],
  providers: [],
  exports: [
    FormErrorsComponent,
    AlertComponent,
    DropdownComponent,
    LoadingOverlayComponent,
    SafeUrlPipe
  ]
})
export class SharedModule {
}
