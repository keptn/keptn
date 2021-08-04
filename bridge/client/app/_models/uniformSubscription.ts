import { DtAutoComplete, DtFilter, DtFilterArray } from './dt-filter';
import { DtFilterFieldChangeEvent } from '@dynatrace/barista-components/filter-field';

export class UniformSubscription {
  public topics: string[] = [];
  public filter!: {
    project: [string],
    stage: [string],
    service: string[]
  };
  public parameters: {key: string, value: string, visible: boolean}[] = [];
  public name!: string;
  public expanded = false;
  private _filter?: DtFilterArray[];

  public static fromJSON(data: unknown): UniformSubscription {
    return Object.assign(new this(), data);
  }

  public get project(): string {
    return this.filter.project[0];
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
      ...this.filter.service.map(service => {
        return [
            data.autocomplete[0],
            {name: service}
          ] as DtFilterArray;
        }
      )
    ];
    if (filter.length !== this._filter?.length) {
      this._filter = filter;
    }
    return this._filter;
  }

  // tslint:disable-next-line:no-any
  public filterChanged(event: DtFilterFieldChangeEvent<any>) { // can't set another type because of "is not assignable to..."
    event = event as DtFilterFieldChangeEvent<DtAutoComplete>;
    this.filter.service = event.filters.reduce((filters: string[], currentFilter: DtAutoComplete[]) => {
      filters.push(currentFilter[1].name);
      return filters;
    }, []);
  }
}
