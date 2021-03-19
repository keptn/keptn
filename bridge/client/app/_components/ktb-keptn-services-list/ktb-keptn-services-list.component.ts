import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  Input,
  OnInit,
  Output,
  ViewChild
} from '@angular/core';
import {DtSort, DtTableDataSource} from '@dynatrace/barista-components/table';
import {KeptnService} from '../../_models/keptn-service';
import {Subscription} from '../../_models/subscription';

@Component({
  selector: 'ktb-keptn-services-list',
  templateUrl: './ktb-keptn-services-list.component.html',
  styleUrls: ['./ktb-keptn-services-list.component.scss']
})
export class KtbKeptnServicesListComponent implements OnInit {

  @ViewChild('sortable', { read: DtSort, static: true }) sortable: DtSort;
  public tableEntries: DtTableDataSource<object> = new DtTableDataSource();
  private _keptnServices: KeptnService[];
  public selectedService: KeptnService;

  @Output() selectedServiceChanged: EventEmitter<KeptnService> = new EventEmitter();

  @Input()
  get keptnServices(): KeptnService[] {
    return this._keptnServices;
  }
  set keptnServices(services: KeptnService[]) {
    if (this._keptnServices !== services) {
      this._keptnServices = services;
      this.tableEntries.data = this._keptnServices;
    }
  }

  constructor() { }

  ngOnInit(): void {
    this.sortable.sort('location', 'desc');
    this.tableEntries.sort = this.sortable;
  }

  public setSelectedService(service: KeptnService) {
    if (this.selectedService !== service) {
      this.selectedService = service;
      this.selectedServiceChanged.emit(service);
    }
  }

  public formatSubscription(subscriptions: Subscription[]): string {
    return subscriptions.reduce((events, subscription) => {
      if (!events.includes(subscription.event)) {
        events.push(subscription.event);
      }
      return events;
    }, []).join(', ');
  }

}
