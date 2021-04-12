import { Component, EventEmitter, OnDestroy, OnInit, Renderer2 } from '@angular/core';
import { Router } from '@angular/router';
import { Claims, SessionService } from '../../../module/core/service/session/session.service';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  templateUrl: './root.component.html',
  styleUrls: ['./root.component.scss']
})
export class RootComponent implements OnInit, OnDestroy {

  loggedInUser: Claims;

  darkTheme: boolean = false;

  destory$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private renderer: Renderer2, private router: Router, private session: SessionService) {
    session.onTokenChange.pipe(takeUntil(this.destory$)).subscribe((token) => {
      if (token) {
        this.loggedInUser = this.session.getClaims();
      } else {
        this.loggedInUser = undefined;
      }
    });
  }

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

  toggleDarkmode() {
    this.darkTheme = !this.darkTheme;
    this.updateTheme();
  }

  updateTheme() {
    if (this.darkTheme) {
      this.renderer.removeClass(document.body, 'light-theme');
      this.renderer.addClass(document.body, 'dark-theme');
      localStorage.setItem("theme", "dark");
    } else {
      this.renderer.removeClass(document.body, 'dark-theme');
      this.renderer.addClass(document.body, 'light-theme');
      localStorage.setItem("theme", "light");
    }
  }

  ngOnInit(): void {
    this.darkTheme = localStorage.getItem("theme") === "dark";
    this.updateTheme();
  }

}
