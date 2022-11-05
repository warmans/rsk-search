import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { RskQuotas } from 'src/app/lib/api-client/models';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';

@Injectable({
  providedIn: 'root'
})
export class QuotaService {

  private quotaSubject$: BehaviorSubject<RskQuotas> = new BehaviorSubject<RskQuotas>({});
  public quotas$: Observable<RskQuotas> = this.quotaSubject$.asObservable();

  constructor(apiClient: SearchAPIClient) {
    apiClient.getQuotaSummary().subscribe((res) => {
      this.quotaSubject$.next(res);
    });
  }
}
