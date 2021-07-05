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
  styleUrls: ['./ktb-sequence-list.component.scss']
})
export class KtbSequenceListComponent implements OnInit, OnDestroy {
  public dataSource: DtTableDataSource<Trace | Sequence> = new DtTableDataSource();
  private unsubscribe$: Subject<void> = new Subject<void>();
  private _sequences: Trace[] = [];
  private _remediations: Sequence[] = [];
  private projectName: string;

  @Input() stage: string;
  @Input() shkeptncontext: string;
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
      if (deployments.some(d => d.shkeptncontext === this.shkeptncontext && d.hasStage(this.stage))) {
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
    this.dataSource.data = [...this.remediations || [], ...this.sequences];
  }

  public isRemediation(row: Sequence | Trace): boolean {
    return row instanceof Sequence;
  }

  public getTraceMessage(trace: Trace): string {
    let message = '';
    const finishedEvent = trace.getFinishedEvent();
    if (finishedEvent?.data.message) {
      message = finishedEvent.data.message;
    } else {
      const failedEvent = trace.findTrace(t => t.data.result === ResultTypes.FAILED);
      if (failedEvent) {
        message = `${failedEvent.source} ${!failedEvent.isFinished() && !failedEvent.isChanged() ? 'started' : failedEvent.isChanged() ? 'changed' : failedEvent.isFinished() ? `finished with result ${failedEvent.data.result}` : ''}`;
      }
    }
    return message;
  }

  public getRemediationLink(remediation: Sequence): string[] {
    return ['/', 'project', this.projectName, 'sequence', remediation.shkeptncontext, 'event', remediation.getStage(this.stage).latestEvent.id];
  }

  public getSequenceLink(trace: Trace): string[] {
    return ['/', 'project', this.projectName, 'sequence', trace.shkeptncontext, 'event', trace.id];
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
