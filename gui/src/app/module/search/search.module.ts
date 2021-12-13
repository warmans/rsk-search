import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SearchComponent } from './page/search/search.component';
import { DslSearchComponent } from './component/dsl-search/dsl-search.component';
import { EpisodeComponent } from './page/episode/episode.component';
import { RouterModule } from '@angular/router';
import { SharedModule } from '../shared/shared.module';
import { GlSearchComponent } from './component/gl-search/gl-search.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgSelectModule } from '@ng-select/ng-select';
import { GlSearchFilterComponent } from './component/gl-search-filter/gl-search-filter.component';
import { NgChartsModule } from 'ng2-charts';
import { SearchStatsComponent } from './component/search-stats/search-stats.component';

@NgModule({
  declarations: [
    SearchComponent,
    DslSearchComponent,
    EpisodeComponent,
    GlSearchComponent,
    GlSearchFilterComponent,
    SearchStatsComponent,
  ],
  imports: [
    CommonModule,
    RouterModule,
    SharedModule,
    FormsModule,
    ReactiveFormsModule,
    NgSelectModule,
    NgChartsModule,
  ],
  exports: [DslSearchComponent, GlSearchComponent]
})
export class SearchModule {
}
