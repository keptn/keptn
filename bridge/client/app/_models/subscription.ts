export class Subscription {
  public event!: string;
  public stages: string[] = [];
  public services: string[] = [];
  public parameters: {key: string, value: string, visible: boolean}[] = [];
  public name!: string;
  public expanded = false;
  private filter: any[] = [];

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
  public getFilter(data: any): any {
    const filter = [
      ...this.stages.map(stage => {
          return [
            data.autocomplete[0],
            {name: stage}
          ];
        }
      ),
      ...this.services.map(services => {
          return [
            data.autocomplete[1],
            {name: services}
          ];
        }
      )
    ];
    if (filter.length !== this.filter.length) {
      this.filter = filter;
    }
    return this.filter;
  }
}
