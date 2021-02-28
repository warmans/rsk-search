import { Component, Input, OnInit } from '@angular/core';
import { AbstractControl, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-form-errors',
  templateUrl: './form-errors.component.html',
  styleUrls: ['./form-errors.component.scss']
})
export class FormErrorsComponent implements OnInit {

  @Input()
  form: FormGroup;

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
