import { Timeframe } from './timeframe';

export type TriggerResponse = { keptnContext: string };

export enum TRIGGER_SEQUENCE {
  DELIVERY,
  EVALUATION,
  CUSTOM,
}

export enum TRIGGER_EVALUATION_TIME {
  TIMEFRAME,
  START_END,
}

export type DeliverySequenceFormData = {
  project?: string;
  service?: string;
  stage?: string;
  image?: string;
  tag?: string;
  labels?: string;
  values?: string;
};

export type EvaluationSequenceFormData = {
  project?: string;
  service?: string;
  stage?: string;
  evaluationType: TRIGGER_EVALUATION_TIME;
  timeframe?: Timeframe;
  timeframeStart?: string; // ISO 8601
  startDatetime?: string; // ISO 8601
  endDatetime?: string; // ISO 8601
  labels?: string;
};

export type CustomSequenceFormData = {
  project?: string;
  service?: string;
  stage?: string;
  sequence?: string;
  labels?: string;
};

export type TriggerSequenceData = {
  project: string;
  stage: string;
  service: string;
  labels?: { [key: string]: string };
  configurationChange?: { values: unknown };
  evaluation?: {
    end?: string; // ISO 8601
    start?: string; // ISO 8601
    timeframe?: string; // e.g. 1h5m
  };
};
