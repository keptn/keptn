import { Component, Input, OnDestroy } from '@angular/core';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { DateUtil } from '../../../_utils/date.utils';
import { SequenceState } from '../../../_models/sequenceState';
import { takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { SubSequence } from '../../../../../shared/interfaces/deployment';
import { EVENT_ICONS } from '../../../_models/event-icons';
import { DtIconType } from '@dynatrace/barista-icons';
import { ResultTypes } from '../../../../../shared/models/result-types';
import { SequenceStatus } from '../../../../../shared/interfaces/sequence';

@Component({
  selector: 'ktb-sequence-list',
  templateUrl: './ktb-sequence-list.component.html',
  styleUrls: [],
})
export class KtbSequenceListComponent implements OnDestroy {
  public dataSource: DtTableDataSource<SubSequence | SequenceState> = new DtTableDataSource();
  private unsubscribe$: Subject<void> = new Subject<void>();
  private _sequences: SubSequence[] = [];
  private _remediations: SequenceState[] = [];
  private projectName?: string;
  public ResultTypes = ResultTypes;
  public SequenceState = SequenceStatus;

  @Input() stage?: string;
  @Input() shkeptncontext?: string;
  @Input()
  get sequences(): SubSequence[] {
    return this._sequences;
  }
  set sequences(sequences: SubSequence[]) {
    if (this._sequences !== sequences) {
      this._sequences = sequences;
      this.updateDataSource();
    }
  }
  @Input()
  get remediations(): SequenceState[] {
    return this._remediations;
  }
  set remediations(remediations: SequenceState[]) {
    if (this._remediations !== remediations) {
      this._remediations = remediations;
      this.updateDataSource();
    }
  }
  constructor(public dateUtil: DateUtil, private route: ActivatedRoute) {
    this.route.paramMap.pipe(takeUntil(this.unsubscribe$)).subscribe((params) => {
      this.projectName = params.get('projectName') ?? undefined;
    });
  }

  private updateDataSource(): void {
    this.dataSource.data = [...this.remediations, ...this.sequences];
  }

  public isRemediation(row: SequenceState | SubSequence): SequenceState | null {
    return row instanceof SequenceState ? row : null;
  }

  public isSubsequence(row: SequenceState | SubSequence): SubSequence | null {
    return row instanceof SequenceState ? null : row;
  }

  public getRemediationLink(remediation: SequenceState): string[] {
    const eventId = this.stage ? remediation.getStage(this.stage)?.latestEvent?.id : undefined;
    return this.projectName && this.stage && eventId
      ? ['/', 'project', this.projectName, 'sequence', remediation.shkeptncontext, 'event', eventId]
      : [];
  }

  public getSequenceLink(subSequence: SubSequence): string[] {
    return this.projectName && this.shkeptncontext
      ? ['/', 'project', this.projectName, 'sequence', this.shkeptncontext, 'event', subSequence.id]
      : [];
  }

  public getEventIcon(subSequence: SubSequence): DtIconType {
    return subSequence.state === SequenceStatus.FINISHED
      ? EVENT_ICONS[subSequence.name] ?? EVENT_ICONS.default
      : EVENT_ICONS.approval;
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
