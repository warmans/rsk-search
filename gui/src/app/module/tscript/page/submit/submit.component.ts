import { Component, OnInit } from '@angular/core';
import { RsksearchTscriptChunk } from '../../../../lib/api-client/models';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { ActivatedRoute, Data } from '@angular/router';
import { Title } from '@angular/platform-browser';

@Component({
  selector: 'app-submit',
  templateUrl: './submit.component.html',
  styleUrls: ['./submit.component.scss']
})
export class SubmitComponent implements OnInit {

  chunk: RsksearchTscriptChunk;

  constructor(
    private route: ActivatedRoute,
    private apiClient: SearchAPIClient,
    private titleService: Title
  ) {
    route.paramMap.subscribe((d: Data) => {
      this.apiClient.searchServiceGetTscriptChunk({ id: d.params['id'] }).subscribe(
        (v) => {
          this.chunk = v;
        }
      );
    });
  }

  ngOnInit(): void {

  }

}
