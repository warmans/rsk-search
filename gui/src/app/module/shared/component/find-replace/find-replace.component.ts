import { ChangeDetectionStrategy, Component, EventEmitter, Output } from '@angular/core';
import { UntypedFormControl, UntypedFormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { NgbPopover } from '@ng-bootstrap/ng-bootstrap';

export interface FindReplace {
  find: string;
  replace: string;
}

@Component({
  selector: 'app-find-replace',
  templateUrl: './find-replace.component.html',
  styleUrls: ['./find-replace.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [NgbPopover, ReactiveFormsModule],
})
export class FindReplaceComponent {
  @Output()
  onSubmit: EventEmitter<FindReplace> = new EventEmitter<FindReplace>();

  form: UntypedFormGroup = new UntypedFormGroup({
    find: new UntypedFormControl('', [Validators.required]),
    replace: new UntypedFormControl(),
  });

  constructor() {}

  submit() {
    this.onSubmit.next({
      find: this.form.get('find').value,
      replace: this.form.get('replace').value,
    });
  }
}
