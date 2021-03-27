import { Component, EventEmitter, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';
import { Claims, SessionService } from '../../../module/core/service/session/session.service';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  templateUrl: './root.component.html',
  styleUrls: ['./root.component.scss']
})
export class RootComponent implements OnDestroy {

  loggedInUser: Claims;

  destory$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private router: Router, private session: SessionService) {
    session.onTokenChange.pipe(takeUntil(this.destory$)).subscribe((token) => {
      if (token) {
        this.loggedInUser = this.session.getClaims();
      } else {
        this.loggedInUser = undefined;
      }
    });
  }

  // @HostListener('document:click', ['$event.path'])
  // public onGlobalClick(targetElementPath: Array<any>) {
  //   let elementRefInPath = targetElementPath.find(e => e === this.sideNav.nativeElement);
  //   if (!elementRefInPath) {
  //     this.showSideNav = false;
  //   }
  // }

  executeSearch(query: string) {
    this.router.navigate(['/search'], { queryParams: { q: query } });
  }

  logout() {
    this.session.destroySession();
    this.loggedInUser = undefined;
    this.router.navigate(['/search']);
  }

  ngOnDestroy(): void {
    this.destory$.next(true);
    this.destory$.complete();
  }

}
