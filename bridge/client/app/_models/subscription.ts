import { DtFilter, DtFilterArray } from './dt-filter';

export class Subscription {
  public event!: string;
  public stages: string[] = [];
  public services: string[] = [];
  public parameters: {key: string, value: string, visible: boolean}[] = [];
  public name!: string;
  public expanded = false;
  private filter: DtFilterArray[] = [];

  static fromJSON(data: unknown) {
    return Object.assign(new this(), data);
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
      ...this.stages.map(stage => {
        return [
            data.autocomplete[0],
            {name: stage}
          ] as DtFilterArray;
        }
      ),
      ...this.services.map(services => {
        return [
            data.autocomplete[1],
            {name: services}
          ] as DtFilterArray;
        }
      )
    ];
    if (filter.length !== this.filter.length) {
      this.filter = filter;
    }
    return this.filter;
  }
}
