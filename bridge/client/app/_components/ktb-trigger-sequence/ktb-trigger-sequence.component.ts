import { Component, EventEmitter, Inject, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { DataService } from '../../_services/data.service';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { Timeframe } from '../../_models/timeframe';
import { DtButton } from '@dynatrace/barista-components/button';
import moment from 'moment';
import { ErrorStateMatcher } from '@angular/material/core';
import { FormControl, FormGroupDirective, NgForm } from '@angular/forms';
import {
  CustomSequenceFormData,
  DeliverySequenceFormData,
  EvaluationSequenceFormData,
  TRIGGER_EVALUATION_TIME,
  TRIGGER_SEQUENCE,
  TriggerEvaluationData,
  TriggerResponse,
  TriggerSequenceData,
} from '../../_models/trigger-sequence';
import { Router } from '@angular/router';
import { takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';

export class ShowErrorStateMatcher implements ErrorStateMatcher {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  isErrorState(control: FormControl | null, form: FormGroupDirective | NgForm | null): boolean {
    return true;
  }
}

export class JsonErrorStateMatcher implements ErrorStateMatcher {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  isErrorState(control: FormControl | null, form: FormGroupDirective | NgForm | null): boolean {
    if (control?.value) {
      return !AppUtils.isValidJson(control.value);
    }
    return false;
  }
}

@Component({
  selector: 'ktb-trigger-sequence[projectName]',
  templateUrl: './ktb-trigger-sequence.component.html',
  styleUrls: ['./ktb-trigger-sequence.component.scss'],
})
export class KtbTriggerSequenceComponent implements OnInit, OnDestroy {
  public TRIGGER_SEQUENCE = TRIGGER_SEQUENCE;
  public TRIGGER_EVALUATION_TIME = TRIGGER_EVALUATION_TIME;
  public state: TRIGGER_SEQUENCE | 'ENTRY' = 'ENTRY';
  public sequenceType: TRIGGER_SEQUENCE = TRIGGER_SEQUENCE.DELIVERY;
  public customSequences: string[] | undefined;
  public selectedService: string | undefined;
  public selectedStage: string | undefined;
  public showErrorStateMatcher = new ShowErrorStateMatcher();
  public jsonErrorStateMatcher = new JsonErrorStateMatcher();
  public isLoading = false;
  public isQualityGatesOnly = false;
  private unsubscribe$: Subject<void> = new Subject<void>();

  public deliveryFormData: DeliverySequenceFormData = {
    image: undefined,
    tag: undefined,
    labels: undefined,
    values: undefined,
  };

  public evaluationFormData: EvaluationSequenceFormData = {
    evaluationType: TRIGGER_EVALUATION_TIME.TIMEFRAME,
    timeframe: undefined,
    timeframeStart: undefined,
    startDatetime: undefined,
    endDatetime: undefined,
    labels: undefined,
  };

  public customFormData: CustomSequenceFormData = {
    sequence: undefined,
    labels: undefined,
  };

  @Input() public projectName!: string;
  @Input() public stage: string | undefined;
  @Input() public stages: string[] = [];
  @Input() public services: string[] = [];

  @Output() public formClosed: EventEmitter<void> = new EventEmitter<void>();

  @ViewChild('timeframeStartButton') timeFrameStartButton?: DtButton;

  constructor(
    private dataService: DataService,
    @Inject(POLLING_INTERVAL_MILLIS) private pollingInterval: number,
    private router: Router
  ) {}

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

  public ngOnInit(): void {
    if (this.stage) {
      this.selectedStage = this.stage;
    }

    this.dataService.isQualityGatesOnly.pipe(takeUntil(this.unsubscribe$)).subscribe((isQualityGatesOnly) => {
      this.isQualityGatesOnly = isQualityGatesOnly;
    });

    AppUtils.createTimer(0, this.pollingInterval)
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        if (this.projectName) {
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
    return start !== undefined && end !== undefined && this.checkStartEndValidity(start, end);
  }

  public isValidJSON(jsonString: string | undefined): boolean {
    if (!jsonString) {
      return true;
    }
    return AppUtils.isValidJson(jsonString);
  }

  public checkStartEndValidity(start: string | undefined, end: string | undefined): boolean {
    const startMoment = moment(start);
    const endMoment = moment(end);
    return startMoment.isBefore(endMoment);
  }

  public triggerSequence(): void {
    this.isLoading = true;

    switch (this.sequenceType) {
      case TRIGGER_SEQUENCE.DELIVERY:
        this.triggerDelivery();
        break;
      case TRIGGER_SEQUENCE.EVALUATION:
        this.triggerEvaluation();
        break;
      case TRIGGER_SEQUENCE.CUSTOM:
        this.triggerCustomSequence();
        break;
    }
  }

  private getImageString(image: string, tag: string): string {
    return image.replace(/\s/g, '') + ':' + tag.replace(/\s/g, '');
  }

  private parseTimeframe(timeframe: Timeframe): string {
    let timeframeString = '';
    timeframeString += timeframe.hours ? timeframe.hours + 'h' : '';
    timeframeString += timeframe.minutes ? timeframe.minutes + 'm' : '';
    timeframeString += timeframe.seconds ? timeframe.seconds + 's' : '';
    timeframeString += timeframe.millis ? timeframe.millis + 'ms' : '';
    timeframeString += timeframe.micros ? timeframe.micros + 'us' : '';

    return timeframeString;
  }

  private parseLabels(labels: string): { [key: string]: string } {
    const labelObj: { [key: string]: string } = {};
    const lbls = labels.replace(/\s/g, '').split(',');
    for (const label of lbls) {
      const parts = label.split('=');
      if (parts[1]) {
        labelObj[parts[0]] = parts[1];
      }
    }

    return labelObj;
  }

  private triggerDelivery(): void {
    const data: TriggerSequenceData = {
      project: this.projectName || '',
      stage: this.selectedStage || '',
      service: this.selectedService || '',
    };

    if (this.deliveryFormData.labels && this.deliveryFormData.labels.trim() !== '') {
      data.labels = this.parseLabels(this.deliveryFormData.labels);
    }

    let valuesObj = {};
    if (this.deliveryFormData.values) {
      valuesObj = JSON.parse(this.deliveryFormData.values);
    }
    data.configurationChange = {
      values: {
        ...valuesObj,
        image: this.getImageString(this.deliveryFormData.image || '', this.deliveryFormData.tag || ''),
      },
    };

    this.dataService.triggerDelivery(data).subscribe(
      (res) => {
        this.handleResponse(res);
      },
      () => {
        this.isLoading = false;
      }
    );
  }

  private triggerEvaluation(): void {
    const data: TriggerEvaluationData = {
      project: this.projectName || '',
      stage: this.selectedStage || '',
      service: this.selectedService || '',
      evaluation: {},
    };
    if (this.evaluationFormData.labels && this.evaluationFormData.labels.trim() !== '') {
      data.evaluation.labels = this.parseLabels(this.evaluationFormData.labels);
    }

    if (this.evaluationFormData.evaluationType === TRIGGER_EVALUATION_TIME.TIMEFRAME) {
      if (this.evaluationFormData.timeframe) {
        data.evaluation.timeframe = this.parseTimeframe(this.evaluationFormData.timeframe);
        data.evaluation.start = moment(this.evaluationFormData.timeframeStart || undefined).toISOString();
      }
    } else if (this.evaluationFormData.evaluationType === TRIGGER_EVALUATION_TIME.START_END) {
      data.evaluation.start = moment(this.evaluationFormData.startDatetime || undefined).toISOString();
      data.evaluation.end = moment(this.evaluationFormData.endDatetime || undefined).toISOString();
    }

    this.dataService.triggerEvaluation(data).subscribe(
      (res) => {
        this.handleResponse(res);
      },
      () => {
        this.isLoading = false;
      }
    );
  }

  private triggerCustomSequence(): void {
    const data: TriggerSequenceData = {
      project: this.projectName || '',
      stage: this.selectedStage || '',
      service: this.selectedService || '',
    };

    if (this.customFormData.labels && this.customFormData.labels.trim() !== '') {
      data.labels = this.parseLabels(this.customFormData.labels);
    }

    this.dataService.triggerCustomSequence(data, this.customFormData.sequence || '').subscribe(
      (res) => {
        this.handleResponse(res);
      },
      () => {
        this.isLoading = false;
      }
    );
  }

  private handleResponse(response: TriggerResponse): void {
    let retry = 1;
    AppUtils.createTimer(500, 1000)
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        if (retry === 5) {
          this.unsubscribe$.next();
          this.navigateToSequences(undefined);
        }
        this.dataService.getSequenceByContext(this.projectName || '', response.keptnContext).subscribe(
          (sequence) => {
            if (sequence) {
              this.navigateToSequences(response.keptnContext);
            } else {
              retry++;
            }
          },
          () => {
            // Gracefully fail - and just navigate to sequences
            this.navigateToSequences(undefined);
          }
        );
      });
  }

  private navigateToSequences(keptnContext: string | undefined): void {
    this.isLoading = false;
    this.unsubscribe$.next();

    if (keptnContext) {
      this.router.navigate(['/project', this.projectName, 'sequence', keptnContext, 'stage', this.selectedStage]);
    } else {
      this.router.navigate(['/project', this.projectName]);
    }
  }
}
