import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {DtTableDataSource} from '@dynatrace/barista-components/table';
import {Trace} from '../../_models/trace';
import {DateUtil} from '../../_utils/date.utils';
import {Sequence} from '../../_models/sequence';
import {takeUntil} from 'rxjs/operators';
import {Subject} from 'rxjs';
import {ActivatedRoute} from '@angular/router';
import {ResultTypes} from '../../_models/result-types';
import {DataService} from '../../_services/data.service';

@Component({
  selector: 'ktb-sequence-list',
  templateUrl: './ktb-sequence-list.component.html',
  styleUrls: []
})
export class KtbSequenceListComponent implements OnInit, OnDestroy {
  public dataSource: DtTableDataSource<Trace | Sequence> = new DtTableDataSource();
  private unsubscribe$: Subject<void> = new Subject<void>();
  private _sequences: Trace[] = [];
  private _remediations: Sequence[] = [];
  private projectName?: string;

  @Input() stage?: string;
  @Input() shkeptncontext?: string;
  @Input()
  get sequences() {
    return this._sequences;
  }
  set sequences(sequences: Trace[]) {
    if (this._sequences !== sequences) {
      this._sequences = sequences;
      this._sequences.sort(DateUtil.compareTraceTimesAsc);
      this.updateDataSource();
    }
  }
  @Input()
  get remediations() {
    return this._remediations;
  }
  set remediations(remediations: Sequence[]) {
    if (this._remediations !== remediations) {
      this._remediations = remediations;
      this.updateDataSource();
    }
  }
  constructor(public dateUtil: DateUtil, private route: ActivatedRoute, private dataService: DataService) { }

  ngOnInit(): void {
    this.dataService.changedDeployments.pipe(takeUntil(this.unsubscribe$)).subscribe((deployments) => {
      if (this.stage && deployments.some(d => d.shkeptncontext === this.shkeptncontext && d.hasStage(this.stage as string))) {
        this.updateDataSource();
      }
    });
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.projectName = params.projectName;
      });
  }

  private updateDataSource(): void {
    this.dataSource.data = [...this.remediations, ...this.sequences];
  }

  public isRemediation(row: Sequence | Trace): Sequence | null {
    return row instanceof Sequence ? row : null;
  }

  public isTrace(row: Sequence | Trace): Trace | null {
    return row instanceof Trace ? row : null;
  }

  public toSequence(row: Sequence): Sequence {
    return  row;
  }

  public getTraceMessage(trace: Trace): string {
    let message = '';
    const finishedEvent = trace.getFinishedEvent();
    if (finishedEvent?.data.message) {
      message = finishedEvent.data.message;
    } else {
      const failedEvent = trace.findTrace(t => t.data.result === ResultTypes.FAILED);
      let eventState;

      if (failedEvent) {
        if (!failedEvent.isFinished() && !failedEvent.isChanged()) {
          eventState = 'started';
        } else if (failedEvent.isChanged()) {
          eventState = 'changed';
        } else if (failedEvent.isFinished()) {
          eventState = `finished with result ${failedEvent.data.result}`;
        } else {
          eventState = '';
        }
        message = `${failedEvent.source} ${eventState}`;
      }
    }
    return message;
  }

  public getRemediationLink(remediation: Sequence): string[] {
    const eventId = this.stage ? remediation.getStage(this.stage)?.latestEvent?.id : undefined;
    return this.projectName && this.stage && eventId
      ? ['/', 'project', this.projectName, 'sequence', remediation.shkeptncontext, 'event', eventId]
      : [];
  }

  public getSequenceLink(trace: Trace): string[] {
    return this.projectName ? ['/', 'project', this.projectName, 'sequence', trace.shkeptncontext, 'event', trace.id] : [];
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

}
