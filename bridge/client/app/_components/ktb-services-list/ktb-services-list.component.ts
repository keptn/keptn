import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  HostBinding,
  Input,
  ViewEncapsulation,
} from '@angular/core';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { Service } from '../../_models/service';
import { DateUtil } from '../../_utils/date.utils';
import { DataService } from '../../_services/data.service';
import { Sequence } from '../../_models/sequence';

const DEFAULT_PAGE_SIZE = 3;

@Component({
  selector: 'ktb-services-list',
  templateUrl: './ktb-services-list.component.html',
  styleUrls: ['./ktb-services-list.component.scss'],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbServicesListComponent {
  @HostBinding('class') cls = 'ktb-services-list';
  public ServiceClass = Service;
  public _services: Service[] = [];
  public _pageSize: number = DEFAULT_PAGE_SIZE;
  public dataSource: DtTableDataSource<Service> = new DtTableDataSource<Service>();

  @Input()
  get services(): Service[] {
    return this._services;
  }

  set services(value: Service[]) {
    if (this._services !== value) {
      this._services = value;
      this.updateDataSource();
    }
  }

  get pageSize(): number {
    return this._pageSize;
  }

  set pageSize(value: number) {
    if (this._pageSize !== value) {
      this._pageSize = value;
      this.updateDataSource();
    }
  }

  get DEFAULT_PAGE_SIZE(): number {
    return DEFAULT_PAGE_SIZE;
  }

  constructor(
    public dataService: DataService,
    public dateUtil: DateUtil,
    private _changeDetectorRef: ChangeDetectorRef
  ) {}

  updateDataSource(): void {
    this.services.sort(this.compare());
    this.dataSource = new DtTableDataSource(this.services.slice(0, this.pageSize));
    this._changeDetectorRef.markForCheck();
  }

  private compare() {
    return (a: Service, b: Service): number => {
      if (!a.latestSequence) {
        return 1;
      } else if (!b.latestSequence) {
        return -1;
      } else {
        return new Date(b.latestSequence.time).getTime() - new Date(a.latestSequence.time).getTime();
      }
    };
  }

  toggleAllServices(): void {
    if (this.services.length > this.pageSize) {
      this.pageSize = this.services.length;
    } else if (this.pageSize > DEFAULT_PAGE_SIZE) {
      this.pageSize = DEFAULT_PAGE_SIZE;
    }
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
