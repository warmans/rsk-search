import { Component, EventEmitter, Input, OnDestroy, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { AlertService } from '../../../core/service/alert/alert.service';

@Component({
  selector: 'app-reject-button',
  templateUrl: './reject-button.component.html',
  styleUrls: ['./reject-button.component.scss']
})
export class RejectButtonComponent implements OnInit, OnDestroy {

  @Input()
  contributionId: string;

  @Input()
  rejectAction: (contributionId: string, comment: string) => void = () => {
  };

  modalOpen: boolean = false;

  loading: boolean = false;

  reasonInput: FormControl = new FormControl();

  destroy$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor() {
  }

  ngOnInit(): void {
  }

  ngOnDestroy(): void {
    this.destroy$.next(true);
    this.destroy$.complete();
  }

  showRejectDialog() {
    this.modalOpen = true;
  }

  reject() {
    this.rejectAction(this.contributionId, this.reasonInput.value);
    this.modalOpen = false;
  }

  cancel() {
    this.modalOpen = false;
  }

}
