import { Component, EventEmitter, Input, OnDestroy, OnInit } from '@angular/core';
import { UntypedFormControl, ReactiveFormsModule } from '@angular/forms';
import { NgClass } from '@angular/common';

@Component({
  selector: 'app-reject-button',
  templateUrl: './reject-button.component.html',
  styleUrls: ['./reject-button.component.scss'],
  imports: [NgClass, ReactiveFormsModule],
})
export class RejectButtonComponent implements OnInit, OnDestroy {
  @Input()
  contributionId: string;

  @Input()
  rejectAction: (contributionId: string, comment: string) => void = () => {};

  modalOpen: boolean = false;

  loading: boolean = false;

  reasonInput: UntypedFormControl = new UntypedFormControl();

  destroy$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor() {}

  ngOnInit(): void {}

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
