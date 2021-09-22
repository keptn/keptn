import { ChangeDetectionStrategy, ChangeDetectorRef, Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewEncapsulation } from '@angular/core';
import { DateUtil } from '../../_utils/date.utils';
import { map, switchMap, takeUntil } from 'rxjs/operators';
import { ActivatedRoute } from '@angular/router';
import { DataService } from '../../_services/data.service';
import { Subject } from 'rxjs';
import { Project } from '../../_models/project';
import { Sequence } from '../../_models/sequence';

@Component({
  selector: 'ktb-root-events-list',
  templateUrl: './ktb-root-events-list.component.html',
  styleUrls: ['./ktb-root-events-list.component.scss'],
  host: {
    class: 'ktb-root-events-list',
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbRootEventsListComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  public project?: Project;
  public _events: Sequence[] = [];
  public _selectedEvent?: Sequence;
  public loading = true;

  @Output() readonly selectedEventChange = new EventEmitter<{ sequence: Sequence, stage?: string }>();

  @Input()
  get events(): Sequence[] {
    return this._events;
  }

  set events(value: Sequence[]) {
    if (this._events !== value) {
      this._events = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get selectedEvent(): Sequence | undefined {
    return this._selectedEvent;
  }

  set selectedEvent(value: Sequence | undefined) {
    if (this._selectedEvent !== value) {
      this._selectedEvent = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, public dateUtil: DateUtil,
              private route: ActivatedRoute, private dataService: DataService) {
  }

  ngOnInit(): void {
    this.route.params.pipe(
      map(params => params.projectName),
      switchMap(projectName => this.dataService.getProject(projectName)),
      takeUntil(this.unsubscribe$),
    ).subscribe(project => {
      this.project = project;
    });

    this.dataService.sequencesUpdated.pipe(
      takeUntil(this.unsubscribe$),
    ).subscribe(() => {
      this.loading = false;
      this._changeDetectorRef.markForCheck();
    });
  }

  selectEvent(sequence: Sequence, stage?: string): void {
    this.selectedEvent = sequence;
    this.selectedEventChange.emit({sequence, stage});
  }

  identifyEvent(index: number, item: Sequence): string | undefined {
    return item?.time;
  }

  loadOldSequences(): void {
    if (this.project) {
      this.loading = true;
      this._changeDetectorRef.markForCheck();
      this.dataService.loadOldSequences(this.project);
    }
  }

  stageClick($event: any): void {
    console.log('event', $event);
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

  public getShortType(type: string | undefined): string | undefined {
    return type ? Sequence.getShortType(type) : undefined;
  }
}
