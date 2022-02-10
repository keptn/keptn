import { Component, EventEmitter, Inject, Input, OnInit, Output, ViewChild } from '@angular/core';
import { DataService } from '../../_services/data.service';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { Timeframe } from '../../_models/timeframe';
import { DtButton } from '@dynatrace/barista-components/button';
import moment from 'moment';
import { ErrorStateMatcher } from '@angular/material/core';
import { FormControl, FormGroupDirective, NgForm } from '@angular/forms';

export enum TRIGGER_SEQUENCE {
  DELIVERY,
  EVALUATION,
  CUSTOM,
}

export enum TRIGGER_EVALUATION_TIME {
  TIMEFRAME,
  START_END,
}

export type DeliveryFormData = {
  image: string | undefined;
  tag: string | undefined;
  labels: string | undefined;
  values: string | undefined;
};

export type EvaluationFormData = {
  evaluationType: TRIGGER_EVALUATION_TIME;
  timeframe: Timeframe | undefined;
  timeframeStart: string | undefined; // ISO 8601
  startDatetime: string | undefined; // ISO 8601
  endDatetime: string | undefined; // ISO 8601
  labels: string | undefined;
};

export type CustomFormData = {
  sequence: string | undefined;
  labels: string | undefined;
};

export class ShowErrorStateMatcher implements ErrorStateMatcher {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  isErrorState(control: FormControl | null, form: FormGroupDirective | NgForm | null): boolean {
    return true;
  }
}

@Component({
  selector: 'ktb-trigger-sequence',
  templateUrl: './ktb-trigger-sequence.component.html',
  styleUrls: ['./ktb-trigger-sequence.component.scss'],
})
export class KtbTriggerSequenceComponent implements OnInit {
  public TRIGGER_SEQUENCE = TRIGGER_SEQUENCE;
  public TRIGGER_EVALUATION_TIME = TRIGGER_EVALUATION_TIME;
  public state: TRIGGER_SEQUENCE | 'ENTRY' = 'ENTRY';
  public sequenceType: TRIGGER_SEQUENCE = TRIGGER_SEQUENCE.DELIVERY;
  public services: string[] | undefined;
  public stages: string[] | undefined;
  public customSequences: string[] | undefined;
  public selectedService: string | undefined;
  public selectedStage: string | undefined;
  public showErrorStateMatcher = new ShowErrorStateMatcher();

  public deliveryFormData: DeliveryFormData = {
    image: undefined,
    tag: undefined,
    labels: undefined,
    values: undefined,
  };

  public evaluationFormData: EvaluationFormData = {
    evaluationType: TRIGGER_EVALUATION_TIME.TIMEFRAME,
    timeframe: undefined,
    timeframeStart: undefined,
    startDatetime: undefined,
    endDatetime: undefined,
    labels: undefined,
  };

  public customFormData: CustomFormData = {
    sequence: undefined,
    labels: undefined,
  };

  @Input() public projectName: string | undefined;
  @Input() public stage: string | undefined;
  @Output() public formClosed: EventEmitter<void> = new EventEmitter<void>();

  @ViewChild('timeframeStartButton') timeFrameStartButton?: DtButton;

  constructor(private dataService: DataService, @Inject(POLLING_INTERVAL_MILLIS) private pollingInterval: number) {}

  public ngOnInit(): void {
    if (!this.projectName) {
      throw new Error('Project name is required');
    }

    if (this.stage) {
      this.selectedStage = this.stage;
    }

    AppUtils.createTimer(0, this.pollingInterval).subscribe(() => {
      if (this.projectName) {
        this.dataService.getServiceNames(this.projectName).subscribe((services) => {
          this.services = services;
        });
        this.dataService.getStageNames(this.projectName).subscribe((stages) => {
          this.stages = stages;
        });
        this.dataService.getCustomSequenceNames(this.projectName).subscribe((customSequences) => {
          this.customSequences = customSequences;
        });
      }
    });
  }

  public setFormState(): void {
    this.state = this.sequenceType;
  }

  public isValidString(input: string | undefined): boolean {
    return input !== undefined && input.trim() !== '';
  }

  public isValidTimeframe(timeframe: Timeframe | undefined): boolean {
    return (
      timeframe !== undefined &&
      (timeframe.hours !== undefined ||
        timeframe.minutes !== undefined ||
        timeframe.seconds !== undefined ||
        timeframe.micros !== undefined ||
        timeframe.millis !== undefined)
    );
  }

  public isValidStartEndTime(start: string | undefined, end: string | undefined): boolean {
    if (start === undefined || end === undefined) {
      return false;
    }

    return this.checkStartEndValidity(start, end);
  }

  public checkStartEndValidity(start: string | undefined, end: string | undefined): boolean {
    const startMoment = moment(start);
    const endMoment = moment(end);
    if (startMoment.isAfter(endMoment)) {
      return false;
    }

    return !endMoment.isBefore(startMoment);
  }
}
