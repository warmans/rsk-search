import { Component, EventEmitter, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';

export interface FindReplace {
  find: string;
  replace: string;
}

@Component({
  selector: 'app-find-replace',
  templateUrl: './find-replace.component.html',
  styleUrls: ['./find-replace.component.scss']
})
export class FindReplaceComponent {

  @Output()
  onSubmit: EventEmitter<FindReplace> = new EventEmitter<FindReplace>();

  form: FormGroup = new FormGroup({
    'find': new FormControl('', [Validators.required]),
    'replace': new FormControl(),
  });

  constructor() {
  }

  submit() {
    this.onSubmit.next({
      find: this.form.get('find').value,
      replace: this.form.get('replace').value,
    });
  }

}
