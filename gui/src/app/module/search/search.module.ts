import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SearchComponent } from './page/search/search.component';
import { EpisodeComponent } from './page/episode/episode.component';
import { RouterModule } from '@angular/router';
import { SharedModule } from '../shared/shared.module';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgChartsModule } from 'ng2-charts';
import { SearchStatsComponent } from './component/search-stats/search-stats.component';
import { MatchedRowPosPipe } from './pipe/match-row-pos.pipe';
import { EpisodeListComponent } from './component/episode-list/episode-list.component';
import { EpisodeSummaryComponent } from './component/episode-summary/episode-summary.component';
import { ChangelogComponent } from './page/changelog/changelog.component';
import { TimecodeAccuracyPipe } from './pipe/timecode-accuracy.pipe';
import { SearchBarHelpComponent } from './component/search-bar-help/search-bar-help.component';
import { SearchBarCompatComponent } from './component/search-bar-compat/search-bar-compat.component';
import { SearchBarSuggestionComponent } from 'src/app/module/search/component/search-bar-suggestion/search-bar-suggestion.component';

@NgModule({
  declarations: [
    SearchComponent,
    EpisodeComponent,
    SearchStatsComponent,
    MatchedRowPosPipe,
    EpisodeListComponent,
    EpisodeSummaryComponent,
    ChangelogComponent,
    TimecodeAccuracyPipe,
    SearchBarHelpComponent,
    SearchBarCompatComponent,
    SearchBarSuggestionComponent
  ],
  imports: [
    CommonModule,
    RouterModule,
    SharedModule,
    FormsModule,
    ReactiveFormsModule,
    NgChartsModule,
  ],
  exports: [SearchBarCompatComponent]
})
export class SearchModule {
}
