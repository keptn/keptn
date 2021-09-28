import { ChartOptions, XAxisOptions, YAxisOptions } from 'highcharts';
import { DtChartOptions } from '@dynatrace/barista-components/chart';

export interface HeatmapOptions extends DtChartOptions {
  xAxis: (Omit<XAxisOptions, 'categories'> & { categories: string[] })[];
  yAxis: (Omit<YAxisOptions, 'categories'> & { categories: string[] })[];
  chart: ChartOptions;
}
