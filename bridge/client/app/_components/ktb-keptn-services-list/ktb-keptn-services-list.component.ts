import {
  Component,
  EventEmitter,
  Input,
  OnInit,
  Output,
  ViewChild
} from '@angular/core';
import {DtSort, DtTableDataSource} from '@dynatrace/barista-components/table';
import {UniformRegistration} from "../../_models/uniform-registration";

@Component({
  selector: 'ktb-keptn-services-list',
  templateUrl: './ktb-keptn-services-list.component.html',
  styleUrls: ['./ktb-keptn-services-list.component.scss']
})
export class KtbKeptnServicesListComponent implements OnInit {

  @ViewChild('sortable', { read: DtSort, static: true }) sortable: DtSort;
  public tableEntries: DtTableDataSource<object> = new DtTableDataSource();
  private _uniformRegistrations: UniformRegistration[];
  public selectedService: UniformRegistration;

  @Output() selectedServiceChanged: EventEmitter<UniformRegistration> = new EventEmitter();

  @Input()
  get uniformRegistrations(): UniformRegistration[] {
    return this._uniformRegistrations;
  }
  set uniformRegistrations(services: UniformRegistration[]) {
    if (this._uniformRegistrations !== services) {
      this._uniformRegistrations = services;
      this.tableEntries.data = this._uniformRegistrations;
    }
  }

  ngOnInit(): void {
    this.sortable.sort('name', 'asc');
    this.tableEntries.sort = this.sortable;
  }

  public setSelectedService(service: UniformRegistration) {
    if (this.selectedService !== service) {
      this.selectedService = service;
      this.selectedServiceChanged.emit(service);
    }
  }

  public formatSubscriptions(subscriptions: string[]): string {
    const formatted = []
    subscriptions.forEach(subscription => {formatted.push(subscription.replace('sh.keptn.', ''))});
    return formatted.join(', ');
  }

  public sortData(sortEvent) {
    const isAscending = sortEvent.direction === 'asc';
    if(this._uniformRegistrations) {
      this._uniformRegistrations.sort((a, b) => {
        switch (sortEvent.active) {
          case 'host': return this.compare(a.metadata.hostname, b.metadata.hostname, isAscending);
          case 'namespace': return this.compare(a.metadata.kubernetesmetadata.namespace, b.metadata.kubernetesmetadata.namespace, isAscending);
          case 'deployment': return this.compare(a.metadata.kubernetesmetadata.deploymentname, b.metadata.kubernetesmetadata.deploymentname, isAscending);
          case 'location': return this.compare(a.metadata.location, b.metadata.location, isAscending);
        }
      });

      this.tableEntries.data = this._uniformRegistrations;
    } else {
      this.tableEntries.data = [];
    }
  }

  private compare(a: number | string, b: number | string, isAsc: boolean): number {
    return (a < b ? -1 : 1) * (isAsc ? 1 : -1);
  }
}
