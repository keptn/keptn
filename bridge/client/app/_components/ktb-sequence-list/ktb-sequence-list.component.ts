import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { DateUtil } from '../../_utils/date.utils';
import { Sequence } from '../../_models/sequence';
import { takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { DataService } from '../../_services/data.service';
import { SubSequence } from '../../../../shared/interfaces/deployment';
import { EVENT_ICONS } from '../../_models/event-icons';
import { DtIconType } from '@dynatrace/barista-icons';

@Component({
  selector: 'ktb-sequence-list',
  templateUrl: './ktb-sequence-list.component.html',
  styleUrls: [],
})
export class KtbSequenceListComponent implements OnInit, OnDestroy {
  public dataSource: DtTableDataSource<SubSequence | Sequence> = new DtTableDataSource();
  private unsubscribe$: Subject<void> = new Subject<void>();
  private _sequences: SubSequence[] = [];
  private _remediations: Sequence[] = [];
  private projectName?: string;

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
  get remediations(): Sequence[] {
    return this._remediations;
  }
  set remediations(remediations: Sequence[]) {
    if (this._remediations !== remediations) {
      this._remediations = remediations;
      this.updateDataSource();
    }
  }
  constructor(public dateUtil: DateUtil, private route: ActivatedRoute, private dataService: DataService) {}

  public ngOnInit(): void {
    this.route.params.pipe(takeUntil(this.unsubscribe$)).subscribe((params) => {
      this.projectName = params.projectName;
    });
  }

  private updateDataSource(): void {
    this.dataSource.data = [...this.remediations, ...this.sequences];
  }

  public isRemediation(row: Sequence | SubSequence): Sequence | null {
    return row instanceof Sequence ? row : null;
  }

  public isSubsequence(row: Sequence | SubSequence): SubSequence | null {
    return row instanceof Sequence ? null : row;
  }

  public getRemediationLink(remediation: Sequence): string[] {
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
    return subSequence.state === 'finished' ? EVENT_ICONS[subSequence.name] : EVENT_ICONS.approval;
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
