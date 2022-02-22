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
  image: string | undefined;
  tag: string | undefined;
  labels: string | undefined;
  values: string | undefined;
};

export type EvaluationSequenceFormData = {
  project?: string;
  service?: string;
  stage?: string;
  evaluationType: TRIGGER_EVALUATION_TIME;
  timeframe: Timeframe | undefined;
  timeframeStart: string | undefined; // ISO 8601
  startDatetime: string | undefined; // ISO 8601
  endDatetime: string | undefined; // ISO 8601
  labels: string | undefined;
};

export type CustomSequenceFormData = {
  project?: string;
  service?: string;
  stage?: string;
  sequence: string | undefined;
  labels: string | undefined;
};

export type TriggerSequenceData = {
  project: string;
  stage: string;
  service: string;
  labels?: { [key: string]: string };
  configurationChange?: { values: unknown };
};

export type TriggerEvaluationData = {
  project: string;
  stage: string;
  service: string;
  evaluation: {
    end?: string; // ISO 8601
    labels?: { [key: string]: string };
    start?: string; // ISO 8601
    timeframe?: string; // e.g. 1h5m
  };
};
