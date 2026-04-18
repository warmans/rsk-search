import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { FormControl, UntypedFormGroup, Validators } from '@angular/forms';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

export interface TranscriptMetadata {
  name: string;
  summary: string;
  release_date: string;
}

@Component({
  selector: 'app-metadata-editor',
  templateUrl: './metadata-editor.component.html',
  styleUrls: ['./metadata-editor.component.scss'],
  standalone: false,
})
export class MetadataEditorComponent implements OnInit, OnDestroy {
  @Input()
  set transcriptMeta(value: TranscriptMetadata) {
    if (!value) {
      return;
    }
    this.form.get('name').setValue(value.name ?? '', { emitEvent: false });
    this.form.get('summary').setValue(value.summary ?? '', { emitEvent: false });
    this.form.get('release_date').setValue(value.release_date ?? '', { emitEvent: false });
  }
  get transcriptMeta(): TranscriptMetadata {
    return {
      summary: this.form.get('summary').value,
      name: this.form.get('name').value,
      release_date: this.form.get('release_date').value,
    };
  }

  @Input()
  allowEdit: boolean = false;

  @Output()
  metadataUpdated: EventEmitter<TranscriptMetadata> = new EventEmitter<TranscriptMetadata>();

  form: UntypedFormGroup = new UntypedFormGroup({
    name: new FormControl<string>(''),
    summary: new FormControl<string>(''),
    release_date: new FormControl<string>('', [Validators.pattern(/^\d{4}-\d{2}-\d{2}$/)]),
  });

  destroy$: Subject<void> = new Subject<void>();

  constructor() {}

  ngOnInit(): void {
    this.form.valueChanges.pipe(takeUntil(this.destroy$)).subscribe((_res) => {
      this.metadataUpdated.next(this.transcriptMeta);
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
