import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, RouterStateSnapshot, UrlTree } from '@angular/router';
import { Observable, of } from 'rxjs';
import { SessionService } from '../core/service/session/session.service';

@Injectable({
  providedIn: 'root',
})
export class CanActivateAdmin {
  constructor(private session: SessionService) {}

  canActivate(_route: ActivatedRouteSnapshot, _state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    return of(this.session.getClaims().approver);
  }
}
