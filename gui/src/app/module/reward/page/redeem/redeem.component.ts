import { Component, EventEmitter, OnInit } from '@angular/core';
import { RskDonationRecipient, RskReward } from '../../../../lib/api-client/models';
import { takeUntil } from 'rxjs/operators';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { AlertService } from '../../../core/service/alert/alert.service';

@Component({
  selector: 'app-redeem',
  templateUrl: './redeem.component.html',
  styleUrls: ['./redeem.component.scss']
})
export class RedeemComponent implements OnInit {

  organizations: RskDonationRecipient[] = [];

  reward: RskReward;

  form: FormGroup = new FormGroup({
    cause: new FormControl('', [Validators.required]),
  });

  loading: boolean[] = [];

  private destroy$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(
    private apiClient: SearchAPIClient,
    private route: ActivatedRoute,
    private titleService: Title,
    private alertService: AlertService,
    private router: Router) {

    titleService.setTitle('Redeem reward');

    route.paramMap.pipe(takeUntil(this.destroy$)).subscribe((d: Data) => {
      if (d.params['id']) {

        this.loading.push(true);
        this.apiClient.listPendingRewards().pipe(takeUntil(this.destroy$)).subscribe((res) => {
          this.reward = res.rewards.find((r) => r.id === d.params['id']);
        }).add(() => this.loading.pop());

        this.loading.push(true);
        this.apiClient.listDonationRecipients({ rewardId: d.params['id'] }).pipe(takeUntil(this.destroy$)).subscribe((res) => {
          this.organizations = res.organizations;
        }).add(() => this.loading.pop());
      }
    });
  }

  ngOnInit(): void {

  }

  ngOnDestroy(): void {
    this.destroy$.next(true);
    this.destroy$.complete();
  }

  submit() {
    this.loading.push(true);
    this.apiClient.claimReward({
      id: this.reward.id,
      body: { id: this.reward.id, donationArgs: { recipient: this.form.get('cause').value } }
    }).pipe(takeUntil(this.destroy$)).subscribe((res) => {
      this.alertService.success('Reward collected sucessfully.');
      this.router.navigate(['/contribute']);
    }).add(() => this.loading.pop());
  }
}
