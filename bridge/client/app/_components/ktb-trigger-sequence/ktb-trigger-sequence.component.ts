import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  Inject,
  Input,
  OnDestroy,
  OnInit,
  Output,
} from '@angular/core';
import { DataService } from '../../_services/data.service';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { Timeframe } from '../../_models/timeframe';
import moment from 'moment';
import { ErrorStateMatcher } from '@angular/material/core';
import { FormControl, FormGroupDirective, NgForm } from '@angular/forms';
import {
  CustomSequenceFormData,
  DeliverySequenceFormData,
  EvaluationSequenceFormData,
  TRIGGER_EVALUATION_TIME,
  TRIGGER_SEQUENCE,
  TriggerResponse,
  TriggerSequenceData,
} from '../../_models/trigger-sequence';
import { Router } from '@angular/router';
import { takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';
import { ICustomSequences } from '../../../../shared/interfaces/custom-sequences';

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
  public DATE_FORMAT = 'YYYY-MM-dd HH:mm:ss';
  public TRIGGER_SEQUENCE = TRIGGER_SEQUENCE;
  public TRIGGER_EVALUATION_TIME = TRIGGER_EVALUATION_TIME;
  public state: TRIGGER_SEQUENCE | 'ENTRY' = 'ENTRY';
  public sequenceType: TRIGGER_SEQUENCE = TRIGGER_SEQUENCE.DELIVERY;
  public customSequences: ICustomSequences | undefined;
  public selectedService: string | undefined;
  public selectedStage: string | undefined;
  public showErrorStateMatcher = new ShowErrorStateMatcher();
  public jsonErrorStateMatcher = new JsonErrorStateMatcher();
  public isLoading = false;
  public isQualityGatesOnly = false;
  public isValidTimeframe = true;
  public isValidStartBeforeEnd = true;
  public isValidStartEndDuration = true;
  private _services: string[] = [];
  private unsubscribe$: Subject<void> = new Subject<void>();

  public deliveryFormData: DeliverySequenceFormData = {};

  public evaluationFormData: EvaluationSequenceFormData = {
    evaluationType: TRIGGER_EVALUATION_TIME.TIMEFRAME,
  };

  public customFormData: CustomSequenceFormData = {};

  @Input() public projectName!: string;
  @Input() public stage: string | undefined;
  @Input() public stages: string[] = [];

  @Input()
  get services(): string[] {
    return this._services;
  }

  set services(services: string[]) {
    if (services) {
      this._services = services;

      if (this.selectedService && !this._services.find((service) => service === this.selectedService)) {
        this.selectedService = undefined;
      }
    }
  }

  @Output() public formClosed: EventEmitter<void> = new EventEmitter<void>();

  constructor(
    private dataService: DataService,
    @Inject(POLLING_INTERVAL_MILLIS) private pollingInterval: number,
    private router: Router,
    private _changeDetectorRef: ChangeDetectorRef
  ) {}

  public ngAfterViewInit(): void {
    // workaround for "Expression has changed after it was checked"
    // for "disable" of the radio-button of #customSequencesSelect
    this._changeDetectorRef.detectChanges();
  }

  public get customSequencesOfStage(): string[] | undefined {
    if (this.selectedStage) {
      return this.customSequences?.[this.selectedStage];
    }
    return undefined;
  }

  public selectedStageChanged(): void {
    if (this.sequenceType === TRIGGER_SEQUENCE.CUSTOM && this.selectedStage && !this.customSequencesOfStage?.length) {
      this.sequenceType = this.isQualityGatesOnly ? TRIGGER_SEQUENCE.EVALUATION : TRIGGER_SEQUENCE.DELIVERY;
    }
    this.customFormData.sequence = undefined;
  }

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
          this.dataService.getCustomSequences(this.projectName).subscribe((customSequences) => {
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

  public isValidStartEndTime(): boolean {
    return (
      this.evaluationFormData.startDatetime !== undefined &&
      this.evaluationFormData.endDatetime !== undefined &&
      this.isValidStartBeforeEnd &&
      this.isValidStartEndDuration
    );
  }

  public isValidJSON(jsonString: string | undefined): boolean {
    if (!jsonString) {
      return true;
    }
    return AppUtils.isValidJson(jsonString);
  }

  private validateStartEndDate(): void {
    if (this.evaluationFormData.startDatetime && this.evaluationFormData.endDatetime) {
      const start = moment(this.evaluationFormData.startDatetime);
      const end = moment(this.evaluationFormData.endDatetime);

      this.isValidStartBeforeEnd = start.isBefore(end);
      this.isValidStartEndDuration = moment.duration(end.diff(start)).asMinutes() >= 1;
    } else {
      this.isValidStartBeforeEnd = true;
      this.isValidStartEndDuration = true;
    }
  }

  public setStartDate(start: string | undefined): void {
    this.evaluationFormData.startDatetime = start;
    this.validateStartEndDate();
  }

  public setEndDate(end: string | undefined): void {
    this.evaluationFormData.endDatetime = end;
    this.validateStartEndDate();
  }

  public setTimeframe(timeframe: Timeframe): void {
    if (!this.isTimeframeEmpty(timeframe)) {
      this.isValidTimeframe =
        (timeframe.hours ?? 0) * 60 +
          (timeframe.minutes ?? 0) +
          (timeframe.seconds ?? 0) / 60 +
          (timeframe.millis ?? 0) / 60_000 +
          (timeframe.micros ?? 0) / 60_000_000 >=
        1;
    } else {
      this.isValidTimeframe = true;
    }

    this.evaluationFormData.timeframe = timeframe;
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
    const data: TriggerSequenceData = {
      project: this.projectName || '',
      stage: this.selectedStage || '',
      service: this.selectedService || '',
    };
    data.evaluation = {};
    if (this.evaluationFormData.labels && this.evaluationFormData.labels.trim() !== '') {
      data.labels = this.parseLabels(this.evaluationFormData.labels);
    }

    if (this.evaluationFormData.evaluationType === TRIGGER_EVALUATION_TIME.TIMEFRAME) {
      if (this.evaluationFormData.timeframe && !this.isTimeframeEmpty(this.evaluationFormData.timeframe)) {
        data.evaluation.timeframe = this.parseTimeframe(this.evaluationFormData.timeframe);
      } else {
        data.evaluation.timeframe = '5m';
      }

      if (this.evaluationFormData.timeframeStart) {
        // This has only to be set, if entered by the user. If not, we can just set the timeframe and let
        // lighthouse-service do the calculation
        data.evaluation.start = moment(this.evaluationFormData.timeframeStart).toISOString();
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

  private isTimeframeEmpty(timeframe: Timeframe): boolean {
    return (
      timeframe.hours === undefined &&
      timeframe.minutes === undefined &&
      timeframe.seconds === undefined &&
      timeframe.millis === undefined &&
      timeframe.micros === undefined
    );
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
