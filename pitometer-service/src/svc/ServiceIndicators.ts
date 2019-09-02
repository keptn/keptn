export interface ServiceIndicators {
  indicators: Indicator[];
}

export interface Indicator {
  name: string;
  source: string;
  query: string;
}
