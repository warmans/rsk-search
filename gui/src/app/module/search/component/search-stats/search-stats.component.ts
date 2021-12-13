import { Component, Input, OnInit } from '@angular/core';
import { RskSearchStats } from '../../../../lib/api-client/models';
import { ChartConfiguration } from 'chart.js';

@Component({
  selector: 'app-search-stats',
  templateUrl: './search-stats.component.html',
  styleUrls: ['./search-stats.component.scss']
})
export class SearchStatsComponent implements OnInit {

  @Input()
  set rawStats(value: { [p: string]: RskSearchStats }) {
    this._rawStats = value;
    if (value) {
      this.updateCharts();
    }
  }

  get rawStats(): { [p: string]: RskSearchStats } {
    return this._rawStats;
  }

  private _rawStats: { [key: string]: RskSearchStats };

  public episodeCountData: ChartConfiguration['data'];

  public lineChartOptions: ChartConfiguration['options'] = {
    elements: {
      point: {
        radius: 1,
        backgroundColor: 'transparent'
      },
      line: {
        tension: 0.2,
        fill: false,
      }
    },
    plugins: {
      legend: { display: false },
      tooltip: {
        position: 'nearest',
        displayColors: true,
        intersect: false,
      }
    },
  };

  constructor() {
  }

  ngOnInit(): void {

  }

  updateCharts() {
    this.episodeCountData = {
      datasets: [],
      labels: undefined,
    };

    for (let rawStatsKey in this._rawStats) {
      if (this._rawStats.hasOwnProperty(rawStatsKey)) {

        if (this.episodeCountData.labels === undefined) {
          // just take the first one as they should be identical anyway.
          this.episodeCountData.labels = this._rawStats[rawStatsKey].labels;
        }
        const set = {
          data: this._rawStats[rawStatsKey].values,
          label: `${rawStatsKey ? rawStatsKey : 'other'} mentions`,
          borderColor: this.getLineColor(rawStatsKey),
          backgroundColor: this.getLineColor(rawStatsKey),
          pointBackgroundColor: this.getLineColor(rawStatsKey),
          pointBorderColor: 'transparent',
          interaction: {
            intersect: false
          },
        };
        this.episodeCountData.datasets.push(set);
      }
    }
  }

  getLineColor(dataName: string): string | undefined {

    // don't know how to use scss classes for plots - see variables.scss :/
    switch (dataName) {
      case 'ricky':
        return 'rgb(249, 223, 22)';
      case 'steve':
        return 'rgb(63, 247, 39)';
      case 'karl':
        return 'rgb(75, 244, 244)';
      default:
        return undefined;
    }
  }
}
