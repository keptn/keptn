import { DtAutoComplete, DtFilterArray } from './dt-filter';
import { DtFilterFieldChangeEvent } from '@dynatrace/barista-components/filter-field';
import { UniformSubscription as us, UniformSubscriptionFilter } from '../../../shared/interfaces/uniform-subscription';
import { DtFilterFieldDefaultDataSourceAutocomplete } from '@dynatrace/barista-components/filter-field/src/filter-field-default-data-source';
import { EventTypes } from '../../../shared/interfaces/event-types';

export class UniformSubscription implements us {
  public id?: string;
  public filter!: UniformSubscriptionFilter;
  public event = '';
  public parameters: { key: string; value: string; visible: boolean }[] = [];
  private _filter?: DtFilterArray[];

  constructor(projectName?: string) {
    this.filter = {
      projects: projectName ? [projectName] : [],
      stages: [],
      services: [],
    };
  }

  public static fromJSON(data: unknown): UniformSubscription {
    return Object.assign(new this(), data);
  }

  public get prefix(): string {
    return this.eventContent.substring(0, this.eventContent.lastIndexOf('.'));
  }

  public get suffix(): string {
    return this.eventContent.split('.').pop() ?? '';
  }

  public get eventContent(): string {
    return this.event.replace(EventTypes.PREFIX, '');
  }

  public get isGlobal(): boolean {
    return !this.filter.projects?.length;
  }

  public get reduced(): Partial<UniformSubscription> {
    const { _filter, ...subscription } = this;
    return subscription;
  }

  public get formattedEvent(): string {
    return this.event.replace('>', '*');
  }

  public setIsGlobal(status: boolean, projectName: string): void {
    if (status) {
      this.filter.projects = [];
    } else if (!this.hasProject(projectName)) {
      if (!this.filter.projects) {
        this.filter.projects = [];
      }
      this.filter.projects.push(projectName);
    }
  }

  public hasProject(projectName: string, includeEmpty = false): boolean {
    return this.filter.projects?.includes(projectName) || (includeEmpty && !this.filter.projects?.length);
  }

  public addParameter(): void {
    this.parameters.push({ key: '', value: '', visible: true });
  }

  public deleteParameter(index: number): void {
    this.parameters.splice(index, 1);
  }

  public getFilter(data?: DtFilterFieldDefaultDataSourceAutocomplete): DtFilterArray[] {
    if (data) {
      const filter = [
        ...(this.filter.stages?.map((stage) => [data.autocomplete[0], { name: stage }] as DtFilterArray) ?? []),
        ...(this.filter.services?.map((service) => [data.autocomplete[1], { name: service }] as DtFilterArray) ?? []),
      ];
      if (filter.length !== this._filter?.length) {
        this._filter = filter;
      }
    } else {
      this._filter = [];
    }
    return this._filter;
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public filterChanged(event: DtFilterFieldChangeEvent<any>, projectName: string): void {
    // can't set another type because of "is not assignable to..."
    const eventCasted = event as DtFilterFieldChangeEvent<DtAutoComplete>;
    const result = eventCasted.filters.reduce(
      (filters: { Stage: string[]; Service: string[] }, filter) => {
        filters[filter[0].name as 'Stage' | 'Service'].push(filter[1].name);
        return filters;
      },
      { Stage: [], Service: [] }
    );
    this.filter.services = result.Service;
    this.filter.stages = result.Stage;
    if (this.filter.projects?.length && this.hasFilter()) {
      this.setIsGlobal(false, projectName);
    }
  }

  public formatFilter(key: 'services' | 'stages' | 'projects'): string {
    return this.filter[key]?.join(', ') || 'all';
  }

  public hasFilter(): boolean {
    return !!(this.filter.stages?.length || this.filter.services?.length);
  }
}
