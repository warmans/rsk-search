import { Component, EventEmitter, OnInit } from '@angular/core';
import { RskIncomingDonation, RskIncomingDonationList } from 'src/app/lib/api-client/models';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { takeUntil } from 'rxjs/operators';
import { LoadingOverlayComponent } from '../../../shared/component/loading-overlay/loading-overlay.component';
import { DecimalPipe, CurrencyPipe } from '@angular/common';

@Component({
  selector: 'app-donations',
  templateUrl: './donations.component.html',
  styleUrls: ['./donations.component.scss'],
  imports: [LoadingOverlayComponent, DecimalPipe, CurrencyPipe],
})
export class DonationsComponent implements OnInit {
  private destroy$ = new EventEmitter<boolean>();

  donations: RskIncomingDonation[] = [];

  loading: boolean = false;
  showMoreInfo: boolean;

  constructor(private apiClient: SearchAPIClient) {}

  ngOnInit(): void {
    this.apiClient
      .listIncomingDonations()
      .pipe(takeUntil(this.destroy$))
      .subscribe((res: RskIncomingDonationList) => {
        this.donations = res.donations;
      });
  }

  ngOnDestroy(): void {
    this.destroy$.emit(true);
    this.destroy$.complete();
  }
}
