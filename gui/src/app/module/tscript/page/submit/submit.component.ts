import { Component, EventEmitter, HostListener, OnDestroy, OnInit, ViewChild } from '@angular/core';
import {
  RsksearchChunkContribution,
  RsksearchContributionState,
  RsksearchTscriptChunk,
} from '../../../../lib/api-client/models';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged, takeUntil } from 'rxjs/operators';
import { SessionService } from '../../../core/service/session/session.service';
import { getFirstOffset, parseTranscript, Tscript } from '../../../shared/lib/tscript';
import { AudioPlayerComponent } from '../../../shared/component/audio-player/audio-player.component';
import { AlertService } from '../../../core/service/alert/alert.service';

@Component({
  selector: 'app-submit',
  templateUrl: './submit.component.html',
  styleUrls: ['./submit.component.scss']
})
export class SubmitComponent implements OnInit, OnDestroy {

  authenticated: boolean = false;
  authError: string;
  chunk: RsksearchTscriptChunk;
  contribution: RsksearchChunkContribution;
  userCanEdit: boolean = true;

  // to stop the caret from getting messed up by updates we need to separate the input
  // data from the output.
  initialTranscript: string = '';
  updatedTranscript: string = '';
  firstOffset: number = -1;

  contentUpdated: Subject<string> = new Subject<string>();

  audioPlayerURL: string = '';
  parsedTscript: Tscript;
  showHelp: boolean = false;
  autoSeek: boolean = localStorage.getItem('pref-autoseek') === null ? true : localStorage.getItem('pref-autoseek') === 'true';

  cStates = RsksearchContributionState;

  loading: boolean[] = [];
  $destroy: EventEmitter<boolean> = new EventEmitter<boolean>();

  @ViewChild('audioPlayer')
  audioPlayer: AudioPlayerComponent;

  @HostListener('document:keydown', ['$event'])
  handleKeyboardEvent(event: KeyboardEvent): boolean {
    if (this.autoSeek) {
      if (event.key === 'Insert') {
        this.audioPlayer.toggle(-3);
        return false;
      }
      if (event.key === 'ScrollLock') {
        this.audioPlayer.play(-3);
        return false;
      }
      if (event.key === 'Pause') {
        this.audioPlayer.play(3);
        return false;
      }
      return true;
    }
  }

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private apiClient: SearchAPIClient,
    private titleService: Title,
    private sessionService: SessionService,
    private alertService: AlertService,
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
          this.setContribution(res);
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
            titleService.setTitle(`contribute :: ${v.id}`);

            this.chunk = v;
            this.audioPlayerURL = `https://storage.googleapis.com/warmans-transcription-audio/${v.id}.mp3`;

            this.setInitialTranscript(this.chunk.raw);
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
      this.updatedTranscript = v;
    });
  }

  setContribution(res: RsksearchChunkContribution) {
    this.titleService.setTitle(`contribute :: ${res.chunkId} :: ${res.id}`);

    this.contribution = res;
    this.setInitialTranscript(res.transcript);

    this.userCanEdit = true;
    if (this.sessionService.getClaims()?.author_id !== res.authorId || res.state !== RsksearchContributionState.STATE_PENDING) {
      this.userCanEdit = false;
    }
  }

  setInitialTranscript(text: string) {
    this.initialTranscript = this.getBackup() ? this.getBackup() : text;
    this.contentUpdated.next(this.initialTranscript);
    this.firstOffset = getFirstOffset(this.initialTranscript);
  }

  setUpdatedTranscript(text: string) {
    this.contentUpdated.next(text);
  }

  updatePreview(content: string) {
    this.parsedTscript = parseTranscript(content);
  }

  backupContent(text: string) {
    if (!this.chunk) {
      return;
    }
    console.log('backup');
    localStorage.setItem(`chunk-backup-${(this.contribution) ? this.contribution.id : this.chunk.id}`, text);
  }

  getBackup(): string {
    return localStorage.getItem(`chunk-backup-${(this.contribution) ? this.contribution.id : this.chunk.id}`);
  }

  resetToRaw() {
    if (confirm('Really reset editor to raw raw transcript?')) {
      this.initialTranscript = this.contribution ? this.contribution.transcript : this.chunk.raw;
      this.updatePreview(this.initialTranscript);
    }
  }

  toggleAutoseek() {
    this.autoSeek = !this.autoSeek;
    localStorage.setItem('pref-autoseek', this.autoSeek ? 'true' : 'false');
  }

  handleOffsetNavigate(offset: number) {
    if (this.autoSeek) {
      if (offset - this.firstOffset >= 0) {
        this.audioPlayer.seek(offset - this.firstOffset);
      }
    }
  }

  submit() {
    if (!this.contribution) {
      this.loading.push(true);
      this.apiClient.searchServiceCreateChunkContribution({
        chunkId: this.chunk.id,
        body: { chunkId: this.chunk.id, transcript: this.updatedTranscript }
      }).pipe(takeUntil(this.$destroy)).subscribe((res: RsksearchChunkContribution) => {
        this.backupContent(''); // clear backup so that the content always matches what was submitted.
        this.alertService.success('SAVED');
        this.router.navigate(['/chunk', this.chunk.id, 'contrib', res.id]);
      }).add(() => this.loading.shift());
    } else {
      this.loading.push(true);
      this.apiClient.searchServiceUpdateChunkContribution({
        chunkId: this.chunk.id,
        contributionId: this.contribution.id,
        body: { chunkId: this.chunk.id, contributionId: this.contribution.id, transcript: this.updatedTranscript }
      }).pipe(takeUntil(this.$destroy)).subscribe((res) => {
        this.setContribution(res);
        this.alertService.success('UPDATED');
      }).add(() => this.loading.shift());
    }
  }

  markComplete() {
    this.loading.push(true);
    this.apiClient.searchServiceRequestChunkContributionState({
      chunkId: this.chunk.id,
      contributionId: this.contribution.id,
      body: {
        chunkId: this.chunk.id,
        contributionId: this.contribution.id,
        requestState: RsksearchContributionState.STATE_REQUEST_APPROVAL
      }
    }).pipe(takeUntil(this.$destroy)).subscribe((res) => {
      this.setContribution(res);
      this.alertService.success('UPDATED');
    }).add(() => this.loading.shift());
  }

  markIncomplete() {
    this.loading.push(true);
    this.apiClient.searchServiceRequestChunkContributionState({
      chunkId: this.chunk.id,
      contributionId: this.contribution.id,
      body: {
        chunkId: this.chunk.id,
        contributionId: this.contribution.id,
        requestState: RsksearchContributionState.STATE_PENDING
      }
    }).pipe(takeUntil(this.$destroy)).subscribe((res) => {
      this.setContribution(res);
      this.alertService.success('UPDATED');
    }).add(() => this.loading.shift());
  }
}
