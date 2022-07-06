import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceViewComponent } from './ktb-sequence-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { BehaviorSubject, firstValueFrom, of, Subject } from 'rxjs';
import { POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { SequenceFilterMock } from '../../_services/_mockData/sequence-filter.mock';
import { SequencesMock } from '../../_services/_mockData/sequences.mock';
import moment from 'moment';
import { DtQuickFilterDefaultDataSourceAutocomplete } from '@dynatrace/barista-components/quick-filter';
import { KtbSequenceViewModule } from './ktb-sequence-view.module';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RouterTestingModule } from '@angular/router/testing';

describe('KtbSequenceViewComponent', () => {
  let component: KtbSequenceViewComponent;
  let fixture: ComponentFixture<KtbSequenceViewComponent>;
  const queryParams: Subject<Record<string, string | string[]>> = new BehaviorSubject({});

  const projectName = 'sockshop';

  const sequenceFilters = [
    [
      {
        name: 'Stage',
        autocomplete: [],
        showInSidebar: false,
      },
      {
        name: 'dev',
        value: 'dev',
      },
    ],
    [
      {
        name: 'Stage',
        autocomplete: [],
        showInSidebar: false,
      },
      {
        name: 'production',
        value: 'production',
      },
    ],
    [
      {
        name: 'Service',
        autocomplete: [],
        showInSidebar: false,
      },
      {
        name: 'carts',
        value: 'carts',
      },
    ],
    [
      {
        name: 'Sequence',
        autocomplete: [],
        showInSidebar: false,
      },
      {
        name: 'delivery',
        value: 'delivery',
      },
    ],
    [
      {
        name: 'Status',
        autocomplete: [],
        showInSidebar: false,
      },
      {
        name: 'Active',
        value: 'started',
      },
    ],
  ];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbSequenceViewModule, BrowserAnimationsModule, RouterTestingModule, HttpClientTestingModule],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({}),
            paramMap: of(convertToParamMap({ projectName })),
            queryParams,
          },
        },
        { provide: POLLING_INTERVAL_MILLIS, useValue: 0 },
        { provide: ApiService, useClass: ApiServiceMock },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceViewComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should not alter service filters if metadata and sequences match', async () => {
    // given
    const state = await firstValueFrom(component.state$);
    // when
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.mapServiceFilters(SequenceFilterMock, state.sequenceInfo?.sequences);

    expect(getServiceFilter().autocomplete).toEqual([
      { name: 'carts-db', value: 'carts-db' },
      { name: 'carts', value: 'carts' },
    ]);
  });

  it('should add a service if it is in a sequence but not in metadata', async () => {
    // given
    const metadata = SequenceFilterMock;
    metadata.services = metadata.services.splice(1, 1); // remove carts-db
    const state = await firstValueFrom(component.state$);
    getServiceFilter().autocomplete = [{ name: 'carts', value: 'carts' }];

    // when
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.mapServiceFilters(metadata, state.sequenceInfo?.sequences);

    // then
    expect(getServiceFilter().autocomplete).toEqual([
      { name: 'carts', value: 'carts' },
      { name: 'carts-db', value: 'carts-db' },
    ]);
  });

  it('should remove a service from filters if not available in metadata anymore', async () => {
    // given
    const metadata = SequenceFilterMock;
    metadata.services.push('helloservice');
    const state = await firstValueFrom(component.state$);
    getServiceFilter().autocomplete.push({ name: 'helloservice', value: 'helloservice' });

    // when
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.mapServiceFilters(metadata, state.sequenceInfo?.sequences);

    // then
    expect(getServiceFilter().autocomplete).toHaveLength(3);
    // As order gets messed up sometimes, it's safer to test each individually
    expect(getServiceFilter().autocomplete).toContainEqual({
      name: 'carts-db',
      value: 'carts-db',
    });
    expect(getServiceFilter().autocomplete).toContainEqual({ name: 'carts', value: 'carts' });
    expect(getServiceFilter().autocomplete).toContainEqual({
      name: 'helloservice',
      value: 'helloservice',
    });

    // when
    metadata.services.pop();
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.mapServiceFilters(metadata, state.sequenceInfo?.sequences);

    // then
    expect(getServiceFilter().autocomplete).toHaveLength(2);
    expect(getServiceFilter().autocomplete).toContainEqual({ name: 'carts', value: 'carts' });
    expect(getServiceFilter().autocomplete).toContainEqual({
      name: 'carts-db',
      value: 'carts-db',
    });
  });

  it('should show a reload button if older than 1 day', () => {
    // given
    const sequence = SequencesMock[0];
    sequence.time = moment().subtract(1, 'days').subtract(1, 'second').toISOString();

    // when
    const showButton = component.showReloadButton(sequence);

    // then
    expect(showButton).toEqual(true);
  });

  it('should not show a reload button if newer than 1 day', () => {
    // given
    const sequence = SequencesMock[0];
    sequence.time = moment().subtract(1, 'hours').add(1, 'second').toISOString();

    // when
    const showButton = component.showReloadButton(sequence);

    // then
    expect(showButton).toEqual(false);
  });

  it('should save filters on click, load sequences and filter properly', () => {
    // given
    const spySaveSequenceFilters = jest.spyOn(component, 'saveSequenceFilters');
    fixture.detectChanges();

    // when
    component.filtersClicked(
      {
        filters: sequenceFilters,
      },
      SequencesMock,
      projectName
    );
    fixture.detectChanges();

    // then
    expect(spySaveSequenceFilters).toHaveBeenCalledWith(
      {
        Stage: ['dev', 'production'],
        Service: ['carts'],
        Sequence: ['delivery'],
        Status: ['started'],
      },
      projectName
    );
    expect(component.filteredSequences).toEqual([SequencesMock[0]]);
  });

  it('should set sequence filters from query params', () => {
    // given
    const spySetSequenceFilters = jest.spyOn(component, 'setSequenceFilters');
    fixture.detectChanges();

    // when
    queryParams.next({
      Stage: ['dev', 'production'],
      Service: 'carts',
      Sequence: 'delivery',
      Status: 'started',
    });
    fixture.detectChanges();

    // then
    expect(spySetSequenceFilters).toHaveBeenCalledWith(
      {
        Stage: ['dev', 'production'],
        Service: ['carts'],
        Sequence: ['delivery'],
        Status: ['started'],
      },
      projectName
    );
    expect(component.filteredSequences).toEqual([SequencesMock[0]]);
  });

  it('should load sequence filters from local storage', () => {
    // given
    const spyLoadSequenceFilters = jest.spyOn(component, 'loadSequenceFilters');
    const spySetSequenceFilters = jest.spyOn(component, 'setSequenceFilters');
    queryParams.next({
      Stage: ['dev', 'production'],
      Service: 'carts',
      Sequence: 'delivery',
      Status: 'started',
    });
    fixture.detectChanges();

    // when
    queryParams.next({});
    fixture.detectChanges();

    // then
    expect(spyLoadSequenceFilters).toHaveBeenCalled();
    expect(spySetSequenceFilters).toHaveBeenCalledWith(
      {
        Stage: ['dev', 'production'],
        Service: ['carts'],
        Sequence: ['delivery'],
        Status: ['started'],
      },
      projectName
    );
    expect(component.filteredSequences).toEqual([SequencesMock[0]]);
  });

  function getServiceFilter(): { autocomplete: { name: string; value: string }[] } {
    return (component._filterDataSource.data as DtQuickFilterDefaultDataSourceAutocomplete)
      .autocomplete[0] as unknown as {
      autocomplete: { name: string; value: string }[];
    };
  }
});
