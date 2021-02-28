import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SearchComponent } from './page/search/search.component';
import { SearchBarComponent } from './component/search-bar/search-bar.component';

@NgModule({
  declarations: [SearchComponent, SearchBarComponent],
  imports: [
    CommonModule
  ]
})
export class SearchModule {
}
