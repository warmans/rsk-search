import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { RsksearchDialog, RsksearchTscriptChunk } from '../../../../lib/api-client/models';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged, takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-submit',
  templateUrl: './submit.component.html',
  styleUrls: ['./submit.component.scss']
})
export class SubmitComponent implements OnInit, OnDestroy {

  chunk: RsksearchTscriptChunk;

  audioPlayerURL: string = '';

  transcriptEdit: string = '';

  contentUpdated: Subject<string> = new Subject<string>();

  dialogPreview: RsksearchDialog[] = [];

  showHelp: boolean = false;

  $destroy: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private apiClient: SearchAPIClient,
    private titleService: Title
  ) {
    titleService.setTitle('contribute');
    route.paramMap.pipe(takeUntil(this.$destroy)).subscribe((d: Data) => {
      this.apiClient.searchServiceGetTscriptChunk({ id: d.params['id'] }).pipe(takeUntil(this.$destroy)).subscribe(
        (v) => {
          if (!v) {
            return;
          }
          this.chunk = v;
          this.audioPlayerURL = `https://storage.googleapis.com/warmans-transcription-audio/${v.id}.mp3`;
          this.transcriptEdit = this.getBackup() ? this.getBackup() : this.chunk.raw;
          titleService.setTitle(`contribute :: ${v.id}`);
          this.contentUpdated.next(this.transcriptEdit);
        }
      );
    });
  }

  ngOnDestroy(): void {
    this.$destroy.next();
    this.$destroy.complete();
  }

  ngOnInit(): void {
    this.contentUpdated.pipe(takeUntil(this.$destroy), distinctUntilChanged(), debounceTime(100)).subscribe((v) => {

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
    if (confirm('Really reset editor to raw raw transcript?')) {
      this.transcriptEdit = this.chunk.raw;
      this.updatePreview(this.transcriptEdit);
    }
  }

  requestAuth() {
    this.apiClient.searchServiceGetRedditAuthURL().pipe(takeUntil(this.$destroy)).subscribe((res) => {
      document.location.href = res.url;
    });
  }

  submit() {

  }
}
