import * as Highcharts from 'highcharts';
import { SeriesHeatmapDataOptions } from 'highcharts';
import { Trace } from './trace';

export type HeatmapData = (Omit<SeriesHeatmapDataOptions, 'y'> & {y: number, evaluation?: Trace});

export interface HeatmapSeriesOptions extends Highcharts.SeriesHeatmapOptions {
  data: HeatmapData[];
}
