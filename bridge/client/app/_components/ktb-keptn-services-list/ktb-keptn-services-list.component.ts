import {
  Component,
  EventEmitter,
  Input,
  OnInit,
  Output
} from '@angular/core';
import {DtTableDataSource} from '@dynatrace/barista-components/table';
import {UniformRegistration} from "../../_models/uniform-registration";

@Component({
  selector: 'ktb-keptn-services-list',
  templateUrl: './ktb-keptn-services-list.component.html',
  styleUrls: ['./ktb-keptn-services-list.component.scss']
})
export class KtbKeptnServicesListComponent implements OnInit {

  public tableEntries: DtTableDataSource<object> = new DtTableDataSource();
  private _uniformRegistrations: UniformRegistration[];
  public selectedService: UniformRegistration;

  @Output() selectedUniformRegistrationChanged: EventEmitter<UniformRegistration> = new EventEmitter();

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
  }

  public setSelectedService(service: UniformRegistration) {
    if (this.selectedService !== service) {
      this.selectedService = service;
      this.selectedUniformRegistrationChanged.emit(service);
    }
  }

  public formatSubscriptions(subscriptions: string[]): string {
    return subscriptions.join('<br/>');
  }

  public sortData(sortEvent) {
    const isAscending = sortEvent.direction === 'asc';
    if(this._uniformRegistrations) {
      this._uniformRegistrations.sort((a, b) => {
        switch (sortEvent.active) {
          case 'name': return this.compare(a.name, b.name, isAscending);
          case 'host': return (this.compare(a.metadata.hostname, b.metadata.hostname, isAscending) || this.compare(a.name, b.name, true));
          case 'namespace': return this.compare(a.metadata.kubernetesmetadata.namespace, b.metadata.kubernetesmetadata.namespace, isAscending) || this.compare(a.name, b.name, true);
          case 'location': return this.compare(a.metadata.location, b.metadata.location, isAscending) || this.compare(a.name, b.name, true);
        }
      });

      this.tableEntries.data = this._uniformRegistrations;
    } else {
      this.tableEntries.data = [];
    }
  }

  private compare(a: string, b: string, isAsc: boolean): number {
    const result = a.localeCompare(b);
    if (result !== 0 && !isAsc) {
      return -result;
    }
    return result;
  }
}
