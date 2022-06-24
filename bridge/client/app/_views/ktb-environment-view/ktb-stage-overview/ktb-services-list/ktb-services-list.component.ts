import {
  Component,
  DoCheck,
  HostBinding,
  Input,
  IterableDiffer,
  IterableDiffers,
  ViewEncapsulation,
} from '@angular/core';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { Service } from '../../../../_models/service';
import { DateUtil } from '../../../../_utils/date.utils';
import { DataService } from '../../../../_services/data.service';
import { Sequence } from '../../../../_models/sequence';

const DEFAULT_PAGE_SIZE = 3;

@Component({
  selector: 'ktb-services-list',
  templateUrl: './ktb-services-list.component.html',
  styleUrls: [],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
})
export class KtbServicesListComponent implements DoCheck {
  @HostBinding('class') cls = 'ktb-services-list';
  public ServiceClass = Service;
  public _services: Service[] = [];
  public dataSource: DtTableDataSource<Service> = new DtTableDataSource<Service>();
  private _expanded = false;
  private iterableDiffer: IterableDiffer<unknown>;

  @Input()
  get services(): Service[] {
    return this._services;
  }

  set services(value: Service[]) {
    if (this._services !== value) {
      this._services = value;
    }
  }

  get expanded(): boolean {
    return this._expanded;
  }

  set expanded(value: boolean) {
    if (this._expanded !== value) {
      this._expanded = value;
      this.updateDataSource();
    }
  }

  get DEFAULT_PAGE_SIZE(): number {
    return DEFAULT_PAGE_SIZE;
  }

  constructor(public dataService: DataService, public dateUtil: DateUtil, private iterableDiffers: IterableDiffers) {
    this.iterableDiffer = iterableDiffers.find([]).create();
  }

  public ngDoCheck(): void {
    const changes = this.iterableDiffer.diff(this._services);
    if (changes) {
      this.updateDataSource();
    }
  }

  updateDataSource(): void {
    this.dataSource = new DtTableDataSource(this.expanded ? this.services : this.services.slice(0, DEFAULT_PAGE_SIZE));
  }

  toggleAllServices(): void {
    this.expanded = !this.expanded;
  }

  getServiceLink(service: Service): string[] {
    return ['service', service.serviceName, 'context', service.deploymentContext ?? '', 'stage', service.stage];
  }

  getSequenceLink(sequence: Sequence, service: Service): string[] {
    return ['sequence', sequence.shkeptncontext, 'stage', service.stage];
  }

  getImageText(service: Service): string {
    if (!service.deployedImage) {
      return '';
    }

    let text = service.getShortImageName();
    if (service.getImageVersion()) {
      text += ':' + service.getImageVersion();
    }
    return text ?? '';
  }
}
