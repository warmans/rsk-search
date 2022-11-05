import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ImportComponent } from './page/import/import.component';
import { CanActivateAdmin } from './can-activate-admin';
import { ReactiveFormsModule } from '@angular/forms';
import { QuotasComponent } from './page/quotas/quotas.component';

@NgModule({
  declarations: [
    ImportComponent,
    QuotasComponent
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule
  ],
  providers: [
    CanActivateAdmin
  ],
})
export class AdminModule {
}
