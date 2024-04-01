import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {IndexComponent} from "./page/index/index.component";
import {RouterModule} from "@angular/router";
import {SharedModule} from "../shared/shared.module";
import {CatalogWarehouseComponent} from "./component/catalog-warehouse/catalog-warehouse.component";
import {SongSearchComponent} from "./page/song-search/song-search.component";
import {ReactiveFormsModule} from "@angular/forms";

@NgModule({
  declarations: [
    IndexComponent,
    CatalogWarehouseComponent,
    SongSearchComponent,
  ],
  imports: [
    CommonModule,
    RouterModule,
    SharedModule,
    ReactiveFormsModule,
  ],
  exports: [
    IndexComponent
  ]
})
export class MoreShiteModule {
}
