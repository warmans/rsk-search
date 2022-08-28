import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { HttpParams } from '@angular/common/http';
import { FormControl, FormGroup } from '@angular/forms';
import { DomSanitizer, SafeUrl } from '@angular/platform-browser';
import { debounceTime, takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-preview-modal',
  templateUrl: './preview-modal.component.html',
  styleUrls: ['./preview-modal.component.scss']
})
export class PreviewModalComponent implements OnInit, OnDestroy {

  @Input()
  set epid(value: string) {
    this._epid = value;
    this.updateIframeURL();
  }

  get epid(): string {
    return this._epid;
  }

  private _epid: string;

  @Input()
  set startPos(value: number) {
    this.sharePropertiesForm.get('startPos').setValue(value);
    this.updateIframeURL();
  }

  get startPos(): number {
    return this.sharePropertiesForm.get('startPos').value;
  }

  @Input()
  set endPos(value: number) {
    this.sharePropertiesForm.get('endPos').setValue(value);
    this.updateIframeURL();
  }

  get endPos(): number {
    return this.sharePropertiesForm.get('endPos').value;
  }

  @Output()
  closed: EventEmitter<boolean> = new EventEmitter();

  sharePropertiesForm = new FormGroup({
    startPos: new FormControl(),
    endPos: new FormControl(),
  });

  embedURL: string;
  iframeURL: SafeUrl;
  embedCode: string;

  unsubscribe$: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private sanitizer: DomSanitizer) {
  }

  ngOnInit(): void {
    this.sharePropertiesForm.valueChanges.pipe(takeUntil(this.unsubscribe$), debounceTime(1000)).subscribe((v) => {
      this.updateIframeURL();
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next(true);
    this.unsubscribe$.complete();
  }

  updateIframeURL() {
    let params = new HttpParams();
    if (this.epid) {
      params = params.append('epid', this.epid);
    }
    if (this.startPos) {
      params = params.append('start', this.startPos);
    }
    if (this.startPos) {
      params = params.append('end', this.endPos);
    }

    this.embedURL = `${window.location.protocol}//${window.location.host}/embed?${params.toString()}`;
    this.iframeURL = this.sanitizer.bypassSecurityTrustResourceUrl(this.embedURL);
    this.embedCode = `<iframe src="${this.embedURL}"></iframe>`;
  }

  close() {
    this.closed.next(true);
  }
}
