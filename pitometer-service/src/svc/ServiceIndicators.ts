export interface ServiceIndicators {
  indicators: Indicator[];
}

export interface Indicator {
  metric: string;
  source: string;
  query: string;
  queryObject: ServiceIndicatorQueryObject[];
}

interface ServiceIndicatorQueryObject {
  key: string;
  value: string;
}
