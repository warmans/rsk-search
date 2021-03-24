import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SearchComponent } from '../module/search/page/search/search.component';
import { EpisodeComponent } from '../module/search/page/episode/episode.component';
import { SubmitComponent } from '../module/tscript/page/submit/submit.component';
import { RandomComponent } from '../module/tscript/page/random/random.component';
import { AuthorContributionsComponent } from '../module/tscript/page/author-contributions/author-contributions.component';
import { ApproveComponent } from '../module/tscript/page/approve/approve.component';


const routes: Routes = [
  {
    path: 'search',
    component: SearchComponent,
  },
  {
    path: 'ep/:id',
    component: EpisodeComponent,
  },
  {
    path: 'chunk/next',
    component: RandomComponent,
  },
  {
    path: 'chunk/:id',
    component: SubmitComponent,
  },
  {
    path: 'chunk/:id/contrib/:contribution_id',
    component: SubmitComponent,
  },
  {
    path: 'me',
    component: AuthorContributionsComponent,
  },
  {
    path: 'approve/:tscript_id',
    component: ApproveComponent,
  },
  { path: '', redirectTo: '/search', pathMatch: 'full' },
];

@NgModule({
  imports: [
    RouterModule.forRoot(
      routes,
      {
        anchorScrolling: 'enabled',
        onSameUrlNavigation: 'reload',
        scrollPositionRestoration: 'enabled'
      },
    )],
  exports: [RouterModule]
})
export class RootRoutingModule {
}
