import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceViewComponent } from './ktb-sequence-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { BehaviorSubject, firstValueFrom, Observable, of, Subject } from 'rxjs';
import { POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { SequenceFilterMock } from '../../_services/_mockData/sequence-filter.mock';
import {
  SequencesMock,
  SequenceWithoutStagesAndEvents,
  SequenceWithStagesAndStartedSequenceMock,
  SequenceWithStartedSequenceMock,
  SequenceWithUnknownEventStage,
} from '../../_services/_mockData/sequences.mock';
import moment from 'moment';
import { DtQuickFilterDefaultDataSourceAutocomplete } from '@dynatrace/barista-components/quick-filter';
import { KtbSequenceViewModule } from './ktb-sequence-view.module';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RouterTestingModule } from '@angular/router/testing';
import { DataService } from '../../_services/data.service';
import { SequenceStatus } from '../../../../shared/interfaces/sequence';
import { SequenceExecutionResult } from '../../../../shared/interfaces/sequence-execution-result';

// debounceTime seems untestable. No matter if fakeAsync and tick is used, it never happens to go inside the subscribe function
jest.mock('rxjs/operators', () => ({
  ...jest.requireActual('rxjs/operators'),
  debounceTime:
    () =>
    <T>(source: Observable<T>): Observable<T> =>
      source,
}));

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
      imports: [
        KtbSequenceViewModule,
        BrowserAnimationsModule,
        RouterTestingModule.withRoutes([
          {
            path: 'project/:projectName/sequence',
            component: KtbSequenceViewComponent,
          },
        ]),
        HttpClientTestingModule,
      ],
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

  it('should select sequence and load traces', async () => {
    // given
    const spySelectSequence = jest.spyOn(component, 'selectSequence');
    const spyLoadTraces = jest.spyOn(component, 'loadTraces');
    const state = await firstValueFrom(component.state$);

    // when
    component.navigateToBlockingSequence(SequencesMock[1], state);

    // then
    expect(spySelectSequence).toHaveBeenCalled();
    expect(spyLoadTraces).toHaveBeenCalled();
  });

  it('should select update sequence if getSequenceExecution returns the sequence itself', async () => {
    // given
    const dataService = TestBed.inject(DataService);
    const spyUpdateSequence = jest.spyOn(dataService, 'updateSequence');
    const state = await firstValueFrom(component.state$);

    // when
    // currentSequence and returned sequence are the same
    component.navigateToBlockingSequence(SequencesMock[0], state);

    // then
    expect(spyUpdateSequence).toHaveBeenCalled();
  });

  it("should call loadUntilRoot if sequence isn't loaded", () => {
    // given
    const dataService = TestBed.inject(DataService);
    const apiService = TestBed.inject(ApiService);
    const spyLoadUntilRoot = jest.spyOn(dataService, 'loadUntilRoot');
    spyLoadUntilRoot.mockReturnValue(of(SequencesMock.slice(0, 25)));
    jest.spyOn(apiService, 'getSequenceExecution').mockReturnValue(
      of({
        sequenceExecutions: [
          {
            scope: {
              keptnContext: '5b8d8f9c-faba-43f2-9c87-fc0b69d6fc3e',
              stage: 'dev',
            },
          },
        ],
      } as SequenceExecutionResult)
    );

    // when
    component.navigateToBlockingSequence(SequencesMock[0]);

    // then
    expect(spyLoadUntilRoot).toHaveBeenCalled();
    expect(component.currentSequence?.shkeptncontext).toBe('5b8d8f9c-faba-43f2-9c87-fc0b69d6fc3e');
  });

  it("should not run data loading or sequence selection if sequence-execution doesn't return anything", () => {
    // given
    const apiService = TestBed.inject(ApiService);
    jest.spyOn(apiService, 'getSequenceExecution').mockReturnValue(
      of({
        sequenceExecutions: [],
      })
    );

    const dataService = TestBed.inject(DataService);
    const spyLoadUntilRoot = jest.spyOn(dataService, 'loadUntilRoot');
    const spySelectSequence = jest.spyOn(component, 'selectSequence');
    const spyLoadTraces = jest.spyOn(component, 'loadTraces');

    // when
    component.navigateToBlockingSequence(SequencesMock[0]);

    // then
    expect(spySelectSequence).not.toHaveBeenCalled();
    expect(spyLoadTraces).not.toHaveBeenCalled();
    expect(spyLoadUntilRoot).not.toHaveBeenCalled();
  });

  it('should call setTraces with correct stage if stage is switched using selectStage', () => {
    // given
    const spyLoadTraces = jest.spyOn(component, 'loadTraces');
    component.selectSequence({ sequence: SequencesMock[0] });

    // when
    component.selectStage('staging');

    // then
    expect(spyLoadTraces).toHaveBeenCalledWith(component.currentSequence, undefined, 'staging');
  });

  it('should go to the latest stage of the latest event if the sequence does not have any stages and the event is not finished', () => {
    // given
    const apiService = TestBed.inject(ApiService);
    const getSequenceExecutionSpy = jest.spyOn(apiService, 'getSequenceExecution');

    // when
    component.navigateToBlockingSequence(SequenceWithStartedSequenceMock);

    // then
    expect(getSequenceExecutionSpy).toHaveBeenCalledWith({
      project: 'sockshop',
      stage: 'production',
      service: 'carts',
      status: SequenceStatus.STARTED,
      pageSize: 1,
    });
  });

  it('should go to the next stage of the sequence', () => {
    // given
    const apiService = TestBed.inject(ApiService);
    const getSequenceExecutionSpy = jest.spyOn(apiService, 'getSequenceExecution');

    // when
    component.navigateToBlockingSequence(SequenceWithStagesAndStartedSequenceMock);

    // then
    expect(getSequenceExecutionSpy).toHaveBeenCalledWith({
      project: 'sockshop',
      stage: 'staging',
      service: 'carts',
      status: SequenceStatus.STARTED,
      pageSize: 1,
    });
  });

  it('should go to the first stage if the sequence does not have any stages nor events', () => {
    // given
    const apiService = TestBed.inject(ApiService);
    const getSequenceExecutionSpy = jest.spyOn(apiService, 'getSequenceExecution');

    // when
    component.navigateToBlockingSequence(SequenceWithoutStagesAndEvents);

    // then
    expect(getSequenceExecutionSpy).toHaveBeenCalledWith({
      project: 'sockshop',
      stage: 'dev',
      service: 'carts',
      status: SequenceStatus.STARTED,
      pageSize: 1,
    });
  });

  it('should go to the first stage if the given stage is not found', () => {
    // given
    const apiService = TestBed.inject(ApiService);
    const getSequenceExecutionSpy = jest.spyOn(apiService, 'getSequenceExecution');

    // when
    component.navigateToBlockingSequence(SequenceWithUnknownEventStage);

    // then
    expect(getSequenceExecutionSpy).toHaveBeenCalledWith({
      project: 'sockshop',
      stage: 'dev',
      service: 'carts',
      status: SequenceStatus.STARTED,
      pageSize: 1,
    });
  });

  function getServiceFilter(): { autocomplete: { name: string; value: string }[] } {
    return (component._filterDataSource.data as DtQuickFilterDefaultDataSourceAutocomplete)
      .autocomplete[0] as unknown as {
      autocomplete: { name: string; value: string }[];
    };
  }
});
