import {ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {DataService} from '../../_services/data.service';
import {takeUntil} from 'rxjs/operators';
import {ActivatedRoute} from '@angular/router';
import {Subject} from 'rxjs';
import {KeptnService} from '../../_models/keptn-service';

@Component({
  selector: 'ktb-uniform-view',
  templateUrl: './ktb-uniform-view.component.html',
  styleUrls: ['./ktb-uniform-view.component.scss']
})
export class KtbUniformViewComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();

  public keptnServices: KeptnService[];
  public selectedService: KeptnService;

  constructor(private dataService: DataService, private route: ActivatedRoute, private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit(): void {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.dataService.getKeptnServices(params.projectName).subscribe(services => {
          this.keptnServices = services;
          this._changeDetectorRef.markForCheck();
        });
        this.dataService.loadTaskNames(params.projectName);
      });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

  selectService(service: KeptnService) {
    this.selectedService = service;
    this._changeDetectorRef.markForCheck();
  }
}
