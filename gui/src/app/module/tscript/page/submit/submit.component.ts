import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { RsksearchChunkContribution, RsksearchDialog, RsksearchTscriptChunk } from '../../../../lib/api-client/models';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged, takeUntil } from 'rxjs/operators';
import { SessionService } from '../../../core/service/session/session.service';

@Component({
  selector: 'app-submit',
  templateUrl: './submit.component.html',
  styleUrls: ['./submit.component.scss']
})
export class SubmitComponent implements OnInit, OnDestroy {

  chunk: RsksearchTscriptChunk;
  contribution: RsksearchChunkContribution;
  userCanEdit: boolean = true;

  audioPlayerURL: string = '';

  transcriptEdit: string = '';

  contentUpdated: Subject<string> = new Subject<string>();

  dialogPreview: RsksearchDialog[] = [];

  showHelp: boolean = false;

  $destroy: EventEmitter<boolean> = new EventEmitter<boolean>();

  authenticated: boolean = false;
  authError: string;

  loading: boolean[] = [];

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private apiClient: SearchAPIClient,
    private titleService: Title,
    private sessionService: SessionService,
  ) {
    titleService.setTitle('contribute');

    route.paramMap.pipe(takeUntil(this.$destroy)).subscribe((d: Data) => {


      if (d.params['contribution_id']) {

        // load content from existing contribution
        this.loading.push(true);
        this.apiClient.searchServiceGetChunkContribution({
          chunkId: d.params['id'],
          contributionId: d.params['contribution_id']
        }).pipe(takeUntil(this.$destroy)).subscribe((res: RsksearchChunkContribution) => {
          this.contribution = res;

          this.transcriptEdit = this.getBackup() ? this.getBackup() : res.transcript;
          titleService.setTitle(`contribute :: ${res.chunkId} :: ${res.id}`);
          this.contentUpdated.next(this.transcriptEdit);

          if (sessionService.getClaims()?.author_id !== res.authorId) {
            this.userCanEdit = false;
          }
          if (res.state !== 'pending') {
            this.userCanEdit = false;
          }
        }).add(() => this.loading.shift());

        this.loading.push(true);
        this.apiClient.searchServiceGetTscriptChunk({ id: d.params['id'] }).pipe(takeUntil(this.$destroy)).subscribe(
          (v) => {
            if (!v) {
              return;
            }
            this.chunk = v;
            this.audioPlayerURL = `https://storage.googleapis.com/warmans-transcription-audio/${v.id}.mp3`;
          }
        ).add(() => this.loading.shift());

      } else {

        // load everything from the chunk

        this.loading.push(true);
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
        ).add(() => this.loading.shift());
      }

      sessionService.onTokenChange.pipe(takeUntil(this.$destroy)).subscribe((token: string): void => {
        if (token != null) {
          this.authenticated = true;
        }
      });
    });
  }

  ngOnDestroy(): void {
    this.$destroy.next();
    this.$destroy.complete();
  }

  ngOnInit(): void {
    this.contentUpdated.pipe(takeUntil(this.$destroy), distinctUntilChanged(), debounceTime(100)).subscribe((v) => {
      this.backupContent(v);
      this.updatePreview(v);
      this.transcriptEdit = v;
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
    localStorage.setItem(`chunk-backup-${(this.contribution) ? this.contribution.id : this.chunk.id}`, text);
  }

  getBackup(): string {
    return localStorage.getItem(`chunk-backup-${(this.contribution) ? this.contribution.id : this.chunk.id}`);
  }

  resetToRaw() {
    if (confirm('Really reset editor to raw raw transcript?')) {
      this.transcriptEdit = this.contribution ? this.contribution.transcript : this.chunk.raw;
      this.updatePreview(this.transcriptEdit);
    }
  }

  submit() {
    if (!this.contribution) {
      this.loading.push(true);
      this.apiClient.searchServiceCreateChunkContribution({
        chunkId: this.chunk.id,
        body: { chunkId: this.chunk.id, transcript: this.transcriptEdit }
      }).pipe(takeUntil(this.$destroy)).subscribe((res: RsksearchChunkContribution) => {
        this.backupContent(''); // clear backup so that the content always matches what was submitted.
        this.router.navigate(['/chunk', this.chunk.id, 'contrib', res.id]);
      }).add(() => this.loading.shift());
    } else {
      this.loading.push(true);
      this.apiClient.searchServiceUpdateChunkContribution({
        chunkId: this.chunk.id,
        contributionId: this.contribution.id,
        body: { chunkId: this.chunk.id, contributionId: this.contribution.id, transcript: this.transcriptEdit }
      }).pipe(takeUntil(this.$destroy)).subscribe((res) => {
        this.handleTranscriptUpdated(res.transcript);
      }).add(() => this.loading.shift());
    }
  }
}
