import { DtAutoComplete, DtFilter, DtFilterArray } from './dt-filter';
import { DtFilterFieldChangeEvent } from '@dynatrace/barista-components/filter-field';

export class UniformSubscription {
  public topics!: string;
  public filter!: {
    projects: string[] | null,
    stages: string[] | null,
    services: string[] | null
  };
  public parameters: {key: string, value: string, visible: boolean}[] = [];
  public name!: string;
  public expanded = false;
  private _filter?: DtFilterArray[];

  public static fromJSON(data: unknown): UniformSubscription {
    return Object.assign(new this(), data);
  }

  public get project(): string | undefined {
    return this.filter.projects?.[0];
  }

  public get stage(): string | undefined {
    return this.filter.stages?.[0];
  }

  public set stage(stage: string | undefined) {
    this.filter.stages = stage ? [stage] : [];
  }

  public addParameter() {
    this.parameters.push({key: '', value: '', visible: true});
  }

  public deleteParameter(index: number) {
    this.parameters.splice(index, 1);
  }

  // tslint:disable-next-line:no-any
  public getFilter(data: any): DtFilterArray[] {
    data = data as DtFilter;
    const filter = [
      ...this.filter.stages?.map(stage => {
        return [
          data.autocomplete[0],
          {name: stage}
        ] as DtFilterArray;
      }) ?? [],
      ...this.filter.services?.map(service => {
        return [
            data.autocomplete[0],
            {name: service}
          ] as DtFilterArray;
        }
      ) ?? []
    ];
    if (filter.length !== this._filter?.length) {
      this._filter = filter;
    }
    return this._filter;
  }

  // tslint:disable-next-line:no-any
  public filterChanged(event: DtFilterFieldChangeEvent<any>) { // can't set another type because of "is not assignable to..."
    event = event as DtFilterFieldChangeEvent<DtAutoComplete>;
    const result = event.filters.reduce((filters: {Stage: string[], Service: string[]}, filter) => {
      filters[filter[0].name as 'Stage' | 'Service'].push(filter[1].name);
      return filters;
    }, {Stage: [], Service: []});
    this.filter.services = result.Service;
    this.filter.stages = result.Stage;
  }

  public formatFilter(key: 'services' | 'stages'): string {
    return this.filter[key]?.join(', ') || 'all';
  }
}
