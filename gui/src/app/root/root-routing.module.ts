import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {SearchComponent} from '../module/search/page/search/search.component';
import {EpisodeComponent} from '../module/search/page/episode/episode.component';
import {RandomComponent} from '../module/contrib/page/random/random.component';
import {AuthorProfile} from '../module/contrib/page/author-profile/author-profile.component';
import {
  EpisodeChunkContributions
} from '../module/contrib/page/episode-chunk-contributions/episode-chunk-contributions.component';
import {ContributeComponent} from '../module/contrib/page/contribute/contribute.component';
import {RedeemComponent} from '../module/reward/page/redeem/redeem.component';
import {TranscriptChangeComponent} from '../module/contrib/page/transcript-change/transcript-change.component';
import {EpisodeChunkSubmit} from '../module/contrib/page/episode-chunk-submit/episode-chunk-submit.component';
import {ChangelogComponent} from '../module/search/page/changelog/changelog.component';
import {ImportComponent} from '../module/admin/page/import/import.component';
import {CanActivateAdmin} from '../module/admin/can-activate-admin';
import {QuotasComponent} from 'src/app/module/admin/page/quotas/quotas.component';
import {IndexComponent} from "../module/more-shite/page/index/index.component";
import {CatalogWarehouseComponent} from "../module/more-shite/component/catalog-warehouse/catalog-warehouse.component";
import {RandomQuoteComponent} from "../module/more-shite/page/random-quote/random-quote.component";
import {SongSearchComponent} from "../module/more-shite/page/song-search/song-search.component";
import {CommunityProjectsComponent} from "../module/more-shite/page/community-projects/community-projects.component";

const routes: Routes = [
  {
    path: 'search',
    component: SearchComponent,
  },
  {
    path: 'changelog',
    component: ChangelogComponent,
  },
  {
    path: 'ep/:id',
    component: EpisodeComponent,
  },
  {
    path: 'ep/:epid/change',
    component: TranscriptChangeComponent,
  },
  {
    path: 'ep/:epid/change/:change_id',
    component: TranscriptChangeComponent,
  },
  {
    path: 'chunk/next',
    component: RandomComponent,
  },
  {
    path: 'chunk/:id',
    component: EpisodeChunkSubmit,
  },
  {
    path: 'chunk/:id/contrib/:contribution_id',
    component: EpisodeChunkSubmit,
  },
  {
    path: 'me',
    component: AuthorProfile,
  },
  {
    path: 'tscript/:tscript_id',
    component: EpisodeChunkContributions,
  },
  {
    path: 'contribute',
    component: ContributeComponent,
  },
  {
    path: 'reward/redeem/:id',
    component: RedeemComponent,
    data: {
      disableRewardPopup: true
    },
  },
  {
    path: 'admin/import',
    component: ImportComponent,
    canActivate: [CanActivateAdmin],
    data: {
      disableRewardPopup: true
    }
  },
  {
    path: 'admin/quotas',
    component: QuotasComponent,
    canActivate: [CanActivateAdmin],
    data: {
      disableRewardPopup: true
    }
  },
  {
    path: 'more-shite',
    component: IndexComponent,
    children: [
      {
        path: "random-quote",
        component: RandomQuoteComponent,
      },
      {
        path: "song-search",
        component: SongSearchComponent,
      },
      {
        path: "catalog-warehouse",
        component: CatalogWarehouseComponent,
      },
      {
        path: "community-projects",
        component: CommunityProjectsComponent,
      },
      {
        path: '',
        redirectTo: '/more-shite/random-quote',
        pathMatch: 'full',
      },
    ],
  },
  {path: '', redirectTo: '/search', pathMatch: 'full'},
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
  exports: [RouterModule],
})
export class RootRoutingModule {
}
