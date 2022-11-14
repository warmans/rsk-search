import { Component, OnDestroy, OnInit } from '@angular/core';
import { QuotaService } from 'src/app/module/core/service/quota/quota.service';
import { takeUntil } from 'rxjs/operators';
import { RskQuotas } from 'src/app/lib/api-client/models';
import { Subject } from 'rxjs';

@Component({
  selector: 'app-quotas',
  templateUrl: './quotas.component.html',
  styleUrls: ['./quotas.component.scss']
})
export class QuotasComponent implements OnInit, OnDestroy {

  private destroy$: Subject<void> = new Subject<void>();

  quotas: RskQuotas;
  bandwidthQuotaUsedPcnt: number;

  constructor(quotaService: QuotaService) {
    quotaService.quotas$.pipe(takeUntil(this.destroy$)).subscribe((res: RskQuotas) => {
      this.quotas = res;
      this.bandwidthQuotaUsedPcnt = (1 - (res.bandwidthRemainingMib / res.bandwidthTotalMib)) * 100;
    });
  }

  ngOnInit(): void {
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

}
