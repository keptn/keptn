import * as Highcharts from 'highcharts';
import { ChartOptions, XAxisOptions, YAxisOptions } from 'highcharts';

export interface HeatmapOptions extends Highcharts.Options {
  xAxis: (Omit<XAxisOptions, 'categories'> & {categories: string[]})[];
  yAxis: (Omit<YAxisOptions, 'categories'> & {categories: string[]})[];
  chart: ChartOptions;
}
