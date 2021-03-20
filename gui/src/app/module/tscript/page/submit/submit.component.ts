import { Component, OnInit } from '@angular/core';
import { RsksearchDialog, RsksearchTscriptChunk } from '../../../../lib/api-client/models';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { ActivatedRoute, Data } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged } from 'rxjs/operators';

@Component({
  selector: 'app-submit',
  templateUrl: './submit.component.html',
  styleUrls: ['./submit.component.scss']
})
export class SubmitComponent implements OnInit {

  chunk: RsksearchTscriptChunk;

  transcriptEdit: string = '';

  contentUpdated: Subject<string> = new Subject<string>();

  dialogPreview: RsksearchDialog[] = [];

  showHelp: boolean = false;

  constructor(
    private route: ActivatedRoute,
    private apiClient: SearchAPIClient,
    private titleService: Title
  ) {
    titleService.setTitle('contribute');
    route.paramMap.subscribe((d: Data) => {
      this.apiClient.searchServiceGetTscriptChunk({ id: d.params['id'] }).subscribe(
        (v) => {
          this.chunk = v;
          this.transcriptEdit = this.getBackup() ? this.getBackup() : this.chunk.raw;
          titleService.setTitle(`contribute :: ${v.id}`);
          this.contentUpdated.next(this.transcriptEdit);
        }
      );
    });
  }

  ngOnInit(): void {
    this.contentUpdated.pipe(distinctUntilChanged(), debounceTime(100)).subscribe((v) => {

      // if the window is closed or something, try and prevent anythign being lost
      this.backupContent(v);

      this.updatePreview(v);
    });
  }

  updatePreview(content: string) {
    this.dialogPreview = [];
    content.split('\n').forEach((line) => {

      if (line.match(/^#OFFSET:.*/g)) {
        return;
      }
      if (line.match(/^#[/]?SYN.*/g)) {
        return;
      }

      const parts = line.split(':');
      if (parts.length < 2) {
        this.dialogPreview.push({ type: 'unknown', content: parts.join(':') });
      } else {
        const actor = parts.shift();
        this.dialogPreview.push({ type: actor == 'song' ? 'song' : 'chat', actor: actor, content: parts.join(':') });
      }
    });
  }

  handleTranscriptUpdated(newContent: string) {
    this.contentUpdated.next(newContent);
  }

  backupContent(text: string) {
    if (!this.chunk) {
      return;
    }
    localStorage.setItem(`chunk-backup-${this.chunk.id}`, text);
  }

  getBackup(): string {
    return localStorage.getItem(`chunk-backup-${this.chunk.id}`);
  }

  resetToRaw() {
    if (confirm("Really reset editor to raw raw transcript?")) {
      this.transcriptEdit = this.chunk.raw;
      this.updatePreview(this.transcriptEdit);
    }
  }

  submit() {

  }
}
