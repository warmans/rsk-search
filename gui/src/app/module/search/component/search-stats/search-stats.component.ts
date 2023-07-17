import {ChangeDetectionStrategy, Component, Input, OnInit} from '@angular/core';
import {RskSearchStats} from '../../../../lib/api-client/models';
import {ChartConfiguration} from 'chart.js';

@Component({
  selector: 'app-search-stats',
  templateUrl: './search-stats.component.html',
  styleUrls: ['./search-stats.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
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

  showMoreStats: boolean = false;

  episodeCountData: ChartConfiguration['data'];

  lineChartOptions: ChartConfiguration['options'] = {
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
      legend: {display: false},
      tooltip: {
        position: 'nearest',
        displayColors: true,
        intersect: false,
      }
    },
  };

  actorCountData: ChartConfiguration['data'];

  barChartOptions: ChartConfiguration['options'] = {
    indexAxis: 'y',
    // Elements options apply to all of the options unless overridden in a dataset
    // In this case, we are setting the border of each horizontal bar to be 2px wide
    elements: {
      bar: {
        borderWidth: 2,
      }
    },
    responsive: true,
    plugins: {
      legend: {
        display: false,
      },
      title: {
        display: false,
      }
    }
  };

  publicationCountData: ChartConfiguration['data'];


  constructor() {
  }

  ngOnInit(): void {

  }

  updateCharts() {
    this.updateEpisodeCountData();
    this.updateActorCountData();
    this.updatePublicationCountData();
  }

  updateEpisodeCountData() {
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

  updateActorCountData() {
    this.actorCountData = {
      datasets: [{
        data: [],
        pointBorderColor: 'transparent',
        backgroundColor: 'rgba(225, 81, 50, 0.7)',
        borderColor: 'transparent',
      }],
      labels: [],
    };
    for (let actor in this._rawStats) {
      if (this._rawStats.hasOwnProperty(actor) && actor) {
        this.actorCountData.labels.push(actor);
        const actorTotal: number = this._rawStats[actor].values.reduce((prev, cur) => prev + cur);
        this.actorCountData.datasets[0].data.push(actorTotal);
      }
    }
  }

  updatePublicationCountData() {
    this.publicationCountData = {
      datasets: [{
        data: [],
        pointBorderColor: 'transparent',
        backgroundColor: 'rgba(225, 81, 50, 0.7)',
        borderColor: 'transparent',
      }],
      labels: [],
    };

    // create a map of publications -> total count
    let publicationCountMap = {};
    for (let actor in this._rawStats) {
      if (this._rawStats.hasOwnProperty(actor)) {
        this._rawStats[actor].labels.forEach((episode, idx: number) => {
          if (publicationCountMap[this.getPublication(episode)] === undefined) {
            publicationCountMap[this.getPublication(episode)] = 0;
          }
          publicationCountMap[this.getPublication(episode)] += this._rawStats[actor].values[idx]
        })
      }
    }

    // convert the map back into a dataset
    for (let publication in publicationCountMap) {
      if (publicationCountMap[publication] > 0) {
        this.publicationCountData.labels.push(publication);
        this.publicationCountData.datasets[0].data.push(publicationCountMap[publication]);
      }
    }
  }

  getPublication(episode: string): string {
    return episode.split("-")[0]
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
