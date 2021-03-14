import { Injectable } from '@angular/core';
import { HttpErrorResponse, HttpEvent, HttpHandler, HttpInterceptor, HttpRequest, } from '@angular/common/http';
import { EMPTY, Observable, throwError as observableThrowError } from 'rxjs';
import { tap } from 'rxjs/operators';
import { AlertService } from '../service/alert/alert.service';

interface ErrorDetails {
  text: string;
  details?: string[];
}

@Injectable()
export class APIErrorInterceptor implements HttpInterceptor {

  constructor(private alerts: AlertService) {
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
          switch (err.status) {
            case 0:
              this.alerts.danger('API: Network error');
              break;
            default:
              const errText = errorTextFromRPCError(err);
              this.alerts.danger(`API: ${errText.text}`, ...(errText.details || []));
          }
          return observableThrowError(err);
        },
      ),
    );
  }
}

function errorTextFromRPCError(err: HttpErrorResponse): ErrorDetails {
  if (!err) {
    return { text: 'Unknown error' };
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
            if (v?.field_violations?.length > 0) {
              for (const f of v.field_violations) {
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
