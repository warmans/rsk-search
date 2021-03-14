import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SearchComponent } from '../module/search/page/search/search.component';
import { EpisodeComponent } from '../module/search/page/episode/episode.component';


const routes: Routes = [
  {
    path: 'search',
    component: SearchComponent,
  },
  {
    path: 'ep/:id',
    component: EpisodeComponent,
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
