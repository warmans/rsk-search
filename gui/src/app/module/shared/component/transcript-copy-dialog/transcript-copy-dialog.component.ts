import { Component, HostListener, Input, OnInit, ViewChild } from '@angular/core';
import { RskSearchResult } from 'src/app/lib/api-client/models';
import { UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { ClipboardService } from 'src/app/module/core/service/clipboard/clipboard.service';

@Component({
  selector: 'app-transcript-copy-dialog',
  templateUrl: './transcript-copy-dialog.component.html',
  styleUrls: ['./transcript-copy-dialog.component.scss']
})
export class TranscriptCopyDialogComponent implements OnInit {

  @Input()
  payload: RskSearchResult;

  @ViewChild('componentRoot')
  componentRootEl: any;

  optionsOpen: boolean = false;

  options: UntypedFormGroup = new UntypedFormGroup({
    markdown: new UntypedFormControl(),
    includeTimestamps: new UntypedFormControl()
  });

  @HostListener('document:click', ['$event'])
  clickOut(event) {
    if (this.componentRootEl.nativeElement.contains(event.target)) {
      return;
    }
    this.optionsOpen = false;
  }

  constructor(private clipboardService: ClipboardService) {
  }

  ngOnInit(): void {
  }

  copyPlain() {
    let content: string[] = [];
    this.payload.dialogs.forEach((dialog) => {
      content.push(...dialog.transcript.map(t => `${t.actor}: ${t.content}`));
    });
    this.clipboardService.copyTextToClipboard(content.join('\n\n'));
    this.optionsOpen = false;
  }

  copyMarkdown() {
    let content: string[] = [];
    this.payload.dialogs.forEach((dialog) => {
      content.push(...dialog.transcript.map(t => t.isMatchedRow ? `> *__${t.actor}:__ ${t.content}*` : `> __${t.actor}:__ ${t.content}`));
    });
    this.clipboardService.copyTextToClipboard(content.join('\n\n'));
    this.optionsOpen = false;
  }
}
