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
      line: {
        tension: 0.5
      }
    },
    scales: {
      x: {},
      'y-axis-0':
        {
          position: 'left',
        },
      'y-axis-1': {
        position: 'right',
        grid: {
          color: 'rgba(255,0,0,0.3)',
        },
        ticks: {
          color: 'red'
        }
      }
    },
  };

  constructor() {
  }

  ngOnInit(): void {

  }

  updateCharts() {
    this.episodeCountData = {
      datasets: [
        {
          data: this._rawStats['xfm_episode_count'].values,
          label: 'Episodes',
          backgroundColor: 'rgba(148,159,177,0.2)',
          borderColor: 'rgba(148,159,177,1)',
          pointBackgroundColor: 'rgba(148,159,177,1)',
          pointBorderColor: '#fff',
          pointHoverBackgroundColor: '#fff',
          pointHoverBorderColor: 'rgba(148,159,177,0.8)',
          fill: 'origin',
        },
      ],
      labels: this._rawStats['xfm_episode_count'].labels
    };
  }

}
