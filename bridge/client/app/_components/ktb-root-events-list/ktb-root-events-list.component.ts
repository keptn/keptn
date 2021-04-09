import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  Input, OnDestroy,
  OnInit, Output,
  ViewEncapsulation
} from '@angular/core';
import {Root} from '../../_models/root';
import {DateUtil} from '../../_utils/date.utils';
import {filter, takeUntil} from 'rxjs/operators';
import {ActivatedRoute} from '@angular/router';
import {DataService} from '../../_services/data.service';
import {Subject} from 'rxjs';
import {Project} from '../../_models/project';

@Component({
  selector: 'ktb-root-events-list',
  templateUrl: './ktb-root-events-list.component.html',
  styleUrls: ['./ktb-root-events-list.component.scss'],
  host: {
    class: 'ktb-root-events-list'
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbRootEventsListComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  public project: Project;
  public _events: Root[] = [];
  public _selectedEvent: Root = null;
  public loading = true;

  @Output() readonly selectedEventChange = new EventEmitter<any>();

  @Input()
  get events(): Root[] {
    return this._events;
  }
  set events(value: Root[]) {
    if (this._events !== value) {
      this._events = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get selectedEvent(): Root {
    return this._selectedEvent;
  }
  set selectedEvent(value: Root) {
    if (this._selectedEvent !== value) {
      this._selectedEvent = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, public dateUtil: DateUtil, private route: ActivatedRoute, private dataService: DataService) { }

  ngOnInit() {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.dataService.getProject(params.projectName).pipe(
          takeUntil(this.unsubscribe$)
        ).subscribe(project => {
          this.project = project;
        });
        this.dataService.roots.pipe(
          takeUntil(this.unsubscribe$),
          filter(roots => !!roots)
        ).subscribe(() => {
          this.loading = false;
          this._changeDetectorRef.markForCheck();
        });
      });
  }

  selectEvent(root: Root, stage?: String) {
    this.selectedEvent = root;
    this.selectedEventChange.emit({ root, stage });
  }

  identifyEvent(index, item) {
    return item?.time;
  }

  loadOldSequences() {
    this.loading = true;
    this._changeDetectorRef.markForCheck();
    this.dataService.loadOldRoots(this.project);
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
