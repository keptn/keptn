import { ChartItem } from '../../../_interfaces/chart';

const labels = { 1: '2022-02-22 09:22', 2: '2022-02-22 12:03' };

const tooltipLabels = {
  0: 'SLO evaluation of test from 2022-02-22 07:01',
  1: '2022-02-22 09:22',
  2: '2022-02-22 12:03',
};

const data: ChartItem[] = [
  {
    type: 'score-bar',
    label: 'Score',
    points: [
      {
        x: 0,
        y: 33,
        color: '#e30505',
      },
      {
        x: 1,
        y: 66,
        color: '#ffaa00',
      },
      {
        x: 2,
        y: 88,
        color: '#518637',
      },
      {
        x: 3,
        y: 50,
        color: '#a97223',
      },
      {
        x: 4,
        y: 2,
        color: '#f30010',
      },
    ],
  },
  {
    type: 'score-line',
    label: 'Score',
    points: [
      {
        x: 0,
        y: 33,
      },
      {
        x: 1,
        y: 66,
      },
      {
        x: 2,
        y: 88,
      },
      {
        x: 3,
        y: 50,
      },
      {
        x: 4,
        y: 2,
      },
    ],
  },
  {
    type: 'metric-line',
    label: 'My custom metric 1 label',
    points: [
      {
        x: 1,
        y: 30,
      },
      {
        x: 2,
        y: 40,
      },
      {
        x: 3,
        y: 144.5,
      },
      {
        x: 4,
        y: 10,
      },
    ],
  },
  {
    type: 'metric-line',
    label: 'Metric 2',
    invisible: true,
    points: [
      {
        x: 1,
        y: 4,
      },
      {
        x: 2,
        y: 5,
      },
      {
        x: 3,
        y: 12,
      },
      {
        x: 4,
        y: 12,
      },
    ],
  },
];

export { labels, tooltipLabels, data };
