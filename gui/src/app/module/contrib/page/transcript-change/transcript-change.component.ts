import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { takeUntil } from 'rxjs/operators';
import { ActivatedRoute, Data, Router } from '@angular/router';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { Title } from '@angular/platform-browser';
import { SessionService } from '../../../core/service/session/session.service';
import { AlertService } from '../../../core/service/alert/alert.service';
import { RskTranscript, RskTranscriptChange } from '../../../../lib/api-client/models';

@Component({
  selector: 'app-transcript-change',
  templateUrl: './transcript-change.component.html',
  styleUrls: ['./transcript-change.component.scss']
})
export class TranscriptChangeComponent implements OnInit, OnDestroy {

  epID: string;

  transcript: RskTranscript;

  change: RskTranscriptChange;

  $destroy: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private apiClient: SearchAPIClient,
    private titleService: Title,
    private sessionService: SessionService,
    private alertService: AlertService,
  ) {
    titleService.setTitle('Contribute');

    route.paramMap.pipe(takeUntil(this.$destroy)).subscribe((d: Data) => {

      this.epID = d.params['epid'];

      this.apiClient.getTranscript({ epid: this.epID, withRaw: true }).pipe(takeUntil(this.$destroy)).subscribe((res: RskTranscript) => {
        this.transcript = res;
      });
    });
  }

  ngOnDestroy(): void {
    this.$destroy.next();
    this.$destroy.complete();
  }

  ngOnInit(): void {


  }

  handleSave(transcript: string): void {

  }

}
