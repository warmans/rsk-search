import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SessionService {

  public onTokenChange: BehaviorSubject<string> = new BehaviorSubject<string>(null);

  private token: string;

  private claims: Claims;

  constructor() {
    const storedToken = localStorage.getItem('token');
    if (storedToken != '') {
      this.registerToken(storedToken);
    }
  }

  registerToken(token: string) {
    this.token = token;
    localStorage.setItem('token', token);
    this.onTokenChange.next(token);
  }

  getToken(): string {
    return this.token;
  }

  getClaims(): Claims {
    if (this.claims) {
      return this.claims;
    }
    if (!this.token) {
      return null;
    }
    // token is header.payload.signature
    const tokenParts = this.token.split('.');
    if (tokenParts.length !== 3) {
      this.destroySession();
      return;
    }

    const claims = JSON.parse(atob(tokenParts[1]));
    if (!claims) {
      return null;
    }
    this.claims = new Claims(claims.author_id, claims.identity as Identity);
    return this.claims;
  }

  destroySession() {
    localStorage.setItem('token', '');
    this.registerToken(null);
  }
}

export class Claims {
  constructor(readonly author_id: string, readonly identity: Identity) {
  }
}

export class Identity {
  id: string;
  name: string;
  icon_img: string;
}
