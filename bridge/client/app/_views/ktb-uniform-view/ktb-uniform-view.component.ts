import {ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute} from '@angular/router';
import {BehaviorSubject, Observable, Subject} from 'rxjs';
import {UniformRegistration} from "../../_models/uniform-registration";
import {UniformRegistrationLog} from "../../_models/uniform-registration-log";
import {switchMap, takeUntil} from "rxjs/operators";

@Component({
  selector: 'ktb-uniform-view',
  templateUrl: './ktb-uniform-view.component.html',
  styleUrls: ['./ktb-uniform-view.component.scss']
})
export class KtbUniformViewComponent implements OnInit, OnDestroy {
  private selectedUniformRegistrationId$ = new Subject<string>();
  private uniformRegistrationLogsSubject = new BehaviorSubject([]);
  private unsubscribe$ = new Subject();

  public selectedUniformRegistration: UniformRegistration;
  public uniformRegistrations$: Observable<UniformRegistration[]>;
  public uniformRegistrationLogs$: Observable<UniformRegistrationLog[]> = this.uniformRegistrationLogsSubject.asObservable();

  public projectName: string

  constructor(private dataService: DataService, private route: ActivatedRoute, private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit(): void {
    this.route.paramMap.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe(map => {
      this.projectName = map.get('projectName');
    });

    this.selectedUniformRegistrationId$.pipe(
      takeUntil(this.unsubscribe$),
      switchMap(uniformRegistrationId => {
        return this.dataService.getUniformRegistrationLogs(uniformRegistrationId);
      })
    ).subscribe((uniformRegLogs) => {
      uniformRegLogs.sort((a, b) => {
        if (a.time.valueOf() > b.time.valueOf()) return -1;
        if (a.time.valueOf() < b.time.valueOf()) return 1;
        return 0;
      });
      this.uniformRegistrationLogsSubject.next(uniformRegLogs);
    });

    this.uniformRegistrations$ = this.dataService.getUniformRegistrations();
  }

  ngOnDestroy() {
    this.unsubscribe$.next();
  }

  selectUniformRegistration(uniformRegistration: UniformRegistration) {
    this.selectedUniformRegistration = uniformRegistration;
    this.selectedUniformRegistrationId$.next(this.selectedUniformRegistration.id);
  }
}
