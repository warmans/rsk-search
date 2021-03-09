import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SearchComponent } from './page/search/search.component';
import { SearchBarComponent } from './component/search-bar/search-bar.component';
import { DslSearchComponent } from './component/dsl-search/dsl-search.component';
import { EpisodeComponent } from './page/episode/episode.component';
import { RouterModule } from '@angular/router';
import { TranscriptComponent } from './component/transcript/transcript.component';
import { SharedModule } from '../shared/shared.module';

@NgModule({
  declarations: [
    SearchComponent,
    SearchBarComponent,
    DslSearchComponent,
    EpisodeComponent,
    TranscriptComponent,
  ],
  imports: [
    CommonModule,
    RouterModule,
    SharedModule
  ],
  exports: [DslSearchComponent]
})
export class SearchModule {
}
