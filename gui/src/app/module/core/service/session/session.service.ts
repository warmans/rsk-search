import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SessionService {

  public onTokenChange: BehaviorSubject<string> = new BehaviorSubject<string>(null);

  private token: string | null;

  private claims: Claims;

  constructor() {
    const storedToken = localStorage.getItem('token');
    if (storedToken != '' && storedToken != 'null') {
      this.registerToken(storedToken);
    }
  }

  registerToken(token: string | null) {
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
    let claims: any;
    try {
      claims = JSON.parse(SessionService.decodeBase64(tokenParts[1]));
    } catch (err) {
      console.error('failed to decode token', err);
      claims = null;
    }
    if (!claims) {
      this.destroySession();
      return;
    }
    this.claims = new Claims(claims.author_id, claims.approver || false, claims.identity as Identity, claims.oauth_provider);
    return this.claims;
  }

  destroySession() {
    this.registerToken(null);
  }

  // atob doesn't support base64 completely (??). Meaning _ will break it. Just replacing it seems to fix it.
  private static decodeBase64(str: string): string {
    return atob(str.replace(/_/g, '/'));
  }
}


export class Claims {
  constructor(readonly author_id: string, readonly approver: boolean, readonly identity: Identity, readonly oauth_provider: string) {
  }
}

export class Identity {
  id: string;
  name: string;
  icon_img: string;
}
