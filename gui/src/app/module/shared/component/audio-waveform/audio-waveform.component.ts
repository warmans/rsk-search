import {AfterViewInit, ChangeDetectorRef, Component, Input, OnDestroy, OnInit} from '@angular/core';
import WaveSurfer from 'wavesurfer.js';
import RegionsPlugin from 'wavesurfer.js/dist/plugins/regions.esm.js';
import {Region} from "wavesurfer.js/dist/plugins/regions";
import {Router} from "@angular/router";
import ZoomPlugin from "wavesurfer.js/dist/plugins/zoom";

const AUDIO_CONTEXT_MS = 2000;

@Component({
  selector: 'app-audio-waveform',
  standalone: false,
  templateUrl: './audio-waveform.component.html',
  styleUrl: './audio-waveform.component.scss'
})
export class AudioWaveformComponent implements OnInit, AfterViewInit, OnDestroy {

  @Input()
  set url(value: string) {
    this._url = value;

  }

  get url(): string {
    return this._url;
  }

  private _url: string;

  @Input()
  startTimestampMs: number;

  @Input()
  endTimestampMs: number;

  @Input()
  episodeDurationMs: number;

  startContext: number = 0;

  endContext: number = 0;

  wave: WaveSurfer = null;

  region: Region = null;

  loading: boolean = false;

  exportURL: string;

  constructor(private cdr: ChangeDetectorRef, private router: Router) {
  }

  ngOnInit(): void {
  }

  ngAfterViewInit(): void {
    this.render()
  }

  ngOnDestroy(): void {
    this.wave.stop();
    this.wave.destroy();
  }

  render(): void {
    if (!this.wave) {
      this.generateWaveform();
    }
    this.cdr.detectChanges();

    Promise.resolve().then(() => this.wave.load(this.audioSegmentUrl()));
  }

  generateWaveform(): void {
    this.loading = true;
    Promise.resolve(null).then(() => {
      this.wave = WaveSurfer.create({
        container: '#waveform',
        waveColor: '#E15132',
        progressColor: '#E87A62FF',
        autoplay: false,
      });

      // Initialize the Zoom plugin
      this.wave.registerPlugin(
        ZoomPlugin.create({
          // the amount of zoom per wheel step, e.g. 0.5 means a 50% magnification per scroll
          scale: 0.5,
          // Optionally, specify the maximum pixels-per-second factor while zooming
          maxZoom: 200,
        }),
      )

      const wsRegions = this.wave.registerPlugin(RegionsPlugin.create())

      wsRegions.on('region-out', (region: Region) => {
        const currentTime: number = this.wave.getCurrentTime();
        const regionStart: number = region.start;
        if (currentTime > regionStart) {
          this.wave?.pause();
        }
      });

      wsRegions.on('region-updated', (region: Region) => {
        const path: string = this._url.split("?")[0];
        this.exportURL = `${path}?${this.getExportQuerystring()}`;
      })

      this.wave.on('decode', () => {
        this.region = wsRegions.addRegion({
          start: this.startContext / 1000,
          end: this.wave.getDuration() - (this.endContext / 1000),
          content: 'Region to download',
          color: 'rgba(0, 255, 0, 0.1)',
          drag: true,
          resize: true,
        });

        const path: string = this._url.split("?")[0];
        this.exportURL = `${path}?${this.getExportQuerystring()}`;
      });

      this.wave.on('ready', () => {
        this.loading = false;
      });
    });
  }


  onPlayPressed() {
    this.region.play();
  }

  onPausePressed() {
    this.wave.pause();
  }

  audioSegmentUrl(): string {
    this.startContext = this.startTimestampMs > AUDIO_CONTEXT_MS ? AUDIO_CONTEXT_MS : AUDIO_CONTEXT_MS - this.startTimestampMs;
    this.endContext = this.episodeDurationMs > this.endTimestampMs + AUDIO_CONTEXT_MS ? AUDIO_CONTEXT_MS : this.episodeDurationMs - this.endTimestampMs;
    return `${this.url}?ts=${this.startTimestampMs - this.startContext}-${this.endTimestampMs + this.endContext}`
  }

  getExportQuerystring(): string {
    const adjustedStartTimeMs = (this.startTimestampMs - this.startContext) + (this.region.start * 1000)
    const adjustedEndTimeMs = (this.endTimestampMs + this.endContext) - ((this.wave.getDuration() - this.region.end) * 1000)
    return `ts=${Math.floor(adjustedStartTimeMs).toFixed(0)}-${Math.ceil(adjustedEndTimeMs).toFixed(0)}`
  }
}
