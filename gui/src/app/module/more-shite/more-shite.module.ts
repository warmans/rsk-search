import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {IndexComponent} from "./page/index/index.component";
import {RouterModule} from "@angular/router";
import {SharedModule} from "../shared/shared.module";
import {CatalogWarehouseComponent} from "./component/catalog-warehouse/catalog-warehouse.component";

@NgModule({
  declarations: [
    IndexComponent,
    CatalogWarehouseComponent
  ],
  imports: [
    CommonModule,
    RouterModule,
    SharedModule,
  ],
  exports: [
    IndexComponent
  ]
})
export class MoreShiteModule {
}
