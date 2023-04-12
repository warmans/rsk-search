import { Component, Input, OnInit } from '@angular/core';
import { AbstractControl, UntypedFormGroup } from '@angular/forms';

@Component({
  selector: 'app-form-errors',
  templateUrl: './form-errors.component.html',
  styleUrls: ['./form-errors.component.scss']
})
export class FormErrorsComponent implements OnInit {

  @Input()
  form: UntypedFormGroup;

  @Input()
  fieldName: string;

  field: AbstractControl;

  constructor() {
  }

  ngOnInit(): void {
    if (!this.form) {
      return undefined;
    }
    this.field = this.form.get(this.fieldName);
  }
}
