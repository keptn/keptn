import { Component, Input } from '@angular/core';
import { DtTableDataSource } from '@dynatrace/barista-components/table';

@Component({
  selector: 'ktb-service-settings-list',
  templateUrl: './ktb-service-settings-list.component.html',
})
export class KtbServiceSettingsListComponent {
  dataSource = new DtTableDataSource<string>();

  @Input()
  projectName = '';

  @Input()
  isLoading = false;

  private _serviceNames: string[] = [];

  @Input()
  get serviceNames(): string[] {
    return this._serviceNames;
  }

  set serviceNames(values: string[] | null) {
    this._serviceNames = values ?? [];
    this.dataSource = new DtTableDataSource<string>(this._serviceNames);
  }
}
