import { Component, EventEmitter, OnInit } from '@angular/core';
import { RskDonationStats, RskIncomingDonation, RskIncomingDonationList, RskRecipientStats } from '../../../../lib/api-client/models';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-donations',
  templateUrl: './donations.component.html',
  styleUrls: ['./donations.component.scss']
})
export class DonationsComponent implements OnInit {

  private destroy$ = new EventEmitter<boolean>();

  donations: RskIncomingDonation[] = [];

  loading: boolean = false;
  showMoreInfo: boolean;

  constructor(private apiClient: SearchAPIClient) {
  }

  ngOnInit(): void {
    this.apiClient.listIncomingDonations().pipe(takeUntil(this.destroy$)).subscribe((res: RskIncomingDonationList) => {
        this.donations = res.donations;
    });
  }

  ngOnDestroy(): void {
    this.destroy$.emit(true);
    this.destroy$.complete();
  }

}
