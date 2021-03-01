import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SearchComponent } from './page/search/search.component';
import { SearchBarComponent } from './component/search-bar/search-bar.component';
import { DslSearchComponent } from './component/dsl-search/dsl-search.component';

@NgModule({
  declarations: [SearchComponent, SearchBarComponent, DslSearchComponent],
  imports: [
    CommonModule
  ]
})
export class SearchModule {
}
