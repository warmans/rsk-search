import { Injectable } from '@angular/core';
import { HttpEvent, HttpHandler, HttpHeaders, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { SessionService } from '../service/session/session.service';
import { Observable } from 'rxjs';

@Injectable()
export class OutgoingTokenInterceptor implements HttpInterceptor {
  constructor(private session: SessionService) {
  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    return next.handle(this.addToken(req));
  }

  private addToken(request: HttpRequest<any>): HttpRequest<any> {

    if (!this.session.getToken()) {
      return request;
    }

    const headers: { [name: string]: string | string[]; } = {};
    for (const key of request.headers.keys()) {
      headers[key] = request.headers.getAll(key);
    }
    headers['Authorization'] = 'Bearer ' + this.session.getToken();

    return request.clone({ headers: new HttpHeaders(headers) });
  }
}
