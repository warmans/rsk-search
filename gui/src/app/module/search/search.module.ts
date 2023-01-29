import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SearchComponent } from './page/search/search.component';
import { EpisodeComponent } from './page/episode/episode.component';
import { RouterModule } from '@angular/router';
import { SharedModule } from '../shared/shared.module';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgSelectModule } from '@ng-select/ng-select';
import { NgChartsModule } from 'ng2-charts';
import { SearchStatsComponent } from './component/search-stats/search-stats.component';
import { MatchedRowPosPipe } from './pipe/match-row-pos.pipe';
import { EpisodeListComponent } from './component/episode-list/episode-list.component';
import { EpisodeSummaryComponent } from './component/episode-summary/episode-summary.component';
import { ChangelogComponent } from './page/changelog/changelog.component';
import { TimecodeAccuracyPipe } from './pipe/timecode-accuracy.pipe';
import { EmbedModule } from '../embed/embed.module';
import { SearchBarComponent } from './component/search-bar/search-bar.component';
import { SearchBarAutocompleteComponent } from './component/search-bar-autocomplete/search-bar-autocomplete.component';
import { SearchBarMentionComponent } from 'src/app/module/search/component/search-bar-mention/search-bar-mention.component';
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
    SearchBarComponent,
    SearchBarAutocompleteComponent,
    SearchBarMentionComponent,
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
    NgSelectModule,
    NgChartsModule,
    EmbedModule,
  ],
    exports: [SearchBarComponent, SearchBarCompatComponent]
})
export class SearchModule {
}
