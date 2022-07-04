import { Injectable } from '@angular/core';
import { HttpErrorResponse, HttpEvent, HttpHandler, HttpInterceptor, HttpRequest, } from '@angular/common/http';
import { EMPTY, Observable, throwError as observableThrowError } from 'rxjs';
import { tap } from 'rxjs/operators';
import { AlertService } from '../service/alert/alert.service';
import { SessionService } from '../service/session/session.service';

interface ErrorDetails {
  text: string;
  details?: string[];
}

@Injectable()
export class APIErrorInterceptor implements HttpInterceptor {

  constructor(private alerts: AlertService, private session: SessionService) {
  }

  intercept(
    req: HttpRequest<any>,
    next: HttpHandler,
  ): Observable<HttpEvent<any>> {
    return next.handle(req).pipe(
      tap(
        (evt: any) => {
        },
        (err: HttpErrorResponse) => {
          if (!err) {
            return EMPTY;
          }
          const errText = errorTextFromRPCError(err);
          switch (err.status) {
            case 401:
              this.session.destroySession();
              this.alerts.danger(`Your session expired or otherwise invalid and has been cleared. Please re-authenticate to access contribution features.`, ...(errText.details || []));
              break;
            default:
              this.alerts.danger(`API Call Failed: ${errText.text}`, ...(errText.details || []));
          }
          return observableThrowError(err);
        },
      ),
    );
  }
}

function errorTextFromRPCError(err: HttpErrorResponse): ErrorDetails {
  if (!err || err.status == 0) {
    return { text: 'Unknown network error' };
  }
  if (err.status === 504) {
    return { text: 'Request timeout' };
  }
  if (typeof err.error === 'object' && err.error && err.error.message) {
    const errDetails: string[] = [];
    if (err?.error?.details?.length > 0) {
      for (const v of err.error.details) {
        switch (v['@type']) {
          case 'type.googleapis.com/google.rpc.BadRequest':
            if (v?.fieldViolations?.length > 0) {
              for (const f of v.fieldViolations) {
                errDetails.push(`Bad request field '${f.field}' (${f.description})`);
              }
            }
            break;
          case 'type.googleapis.com/google.rpc.DebugInfo':
            errDetails.push(v.detail);
            break;
        }
      }
    }
    return { text: `${err.error.message}`, details: errDetails };
  }
  return { text: `${err.status} ${err.statusText === 'OK' ? 'Unknown Error' : err.statusText}` };
}
