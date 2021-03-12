import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SearchComponent } from './page/search/search.component';
import { DslSearchComponent } from './component/dsl-search/dsl-search.component';
import { EpisodeComponent } from './page/episode/episode.component';
import { RouterModule } from '@angular/router';
import { TranscriptComponent } from './component/transcript/transcript.component';
import { SharedModule } from '../shared/shared.module';
import { GlSearchComponent } from './component/gl-search/gl-search.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgSelectModule } from '@ng-select/ng-select';
import { GlSearchFilterComponent } from './component/gl-search-filter/gl-search-filter.component';

@NgModule({
  declarations: [
    SearchComponent,
    DslSearchComponent,
    EpisodeComponent,
    TranscriptComponent,
    GlSearchComponent,
    GlSearchFilterComponent,
  ],
  imports: [
    CommonModule,
    RouterModule,
    SharedModule,
    FormsModule,
    ReactiveFormsModule,
    NgSelectModule
  ],
  exports: [DslSearchComponent, GlSearchComponent]
})
export class SearchModule {
}
