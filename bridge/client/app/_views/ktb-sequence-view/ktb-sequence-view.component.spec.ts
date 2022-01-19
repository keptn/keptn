import { ComponentFixture, TestBed } from '@angular/core/testing';
import { FilterType, KtbSequenceViewComponent } from './ktb-sequence-view.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';
import { POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { SequenceMetadataMock } from '../../_services/_mockData/sequence-metadata.mock';
import { SequencesMock } from '../../_services/_mockData/sequences.mock';
import { ProjectsMock } from '../../_services/_mockData/projects.mock';

describe('KtbEventsListComponent', () => {
  let component: KtbSequenceViewComponent;
  let fixture: ComponentFixture<KtbSequenceViewComponent>;

  const projectName = 'sockshop';
  const activeFilter: FilterType[] = [
    [
      {
        name: 'Service',
        showInSidebar: true,
        autocomplete: [
          {
            name: 'carts-db',
            value: 'carts-db',
          },
          {
            name: 'carts',
            value: 'carts',
          },
        ],
      },
      {
        name: 'carts',
        value: 'carts',
      },
    ],
    [
      {
        name: 'Stage',
        showInSidebar: true,
        autocomplete: [
          {
            name: 'dev',
            value: 'dev',
          },
          {
            name: 'staging',
            value: 'staging',
          },
          {
            name: 'production',
            value: 'production',
          },
        ],
      },
      {
        name: 'production',
        value: 'production',
      },
    ],
  ];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({}),
            params: of({ projectName }),
            queryParams: of({}),
          },
        },
        { provide: POLLING_INTERVAL_MILLIS, useValue: 0 },
        { provide: ApiService, useClass: ApiServiceMock },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should update the latest deployed image', () => {
    // given
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.latestDeployments = SequenceMetadataMock.deployments;
    component.selectedStage = 'staging';
    component.currentSequence = SequencesMock[0];

    // when
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.updateLatestDeployedImage();

    // then
    expect(component.currentLatestDeployedImage).toEqual('carts:0.12.3');
  });

  it('should not alter service filters if metadata and sequences match', () => {
    // given
    /* eslint-disable @typescript-eslint/ban-ts-comment */
    /* @ts-ignore */ // Ignore private property
    component.project = ProjectsMock[0];
    // @ts-ignore // Ignore private property
    component.project.sequences = SequencesMock;
    /* eslint-enable @typescript-eslint/ban-ts-comment */
    component._seqFilters = activeFilter;

    // when
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.mapServiceFilters(SequenceMetadataMock);

    expect(component._seqFilters).toEqual(activeFilter);
  });

  it('should add a service if it is in a sequence but not in metadata', () => {
    // given
    const metadata = SequenceMetadataMock;
    metadata.filter.services = metadata.filter.services.slice(0, 1); // remove carts-db
    /* eslint-disable @typescript-eslint/ban-ts-comment */
    /* @ts-ignore */ // Ignore private property
    component.project = ProjectsMock[0];
    // @ts-ignore // Ignore private property
    component.project.sequences = SequencesMock;
    /* eslint-enable @typescript-eslint/ban-ts-comment */
    const filter = activeFilter;
    filter[0][0].autocomplete = filter[0][0].autocomplete.slice(0, 1); // remove carts-db
    component._seqFilters = filter;

    // when
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.mapServiceFilters(metadata);

    // then
    expect(component._seqFilters).toEqual(activeFilter);
  });

  it('should remove a service from filters if not available in metadata anymore', () => {
    // given
    const metadata = SequenceMetadataMock;
    metadata.filter.services.push('helloservice');
    /* eslint-disable @typescript-eslint/ban-ts-comment */
    /* @ts-ignore */ // Ignore private property
    component.project = ProjectsMock[0];
    // @ts-ignore // Ignore private property
    component.project.sequences = SequencesMock;
    /* eslint-enable @typescript-eslint/ban-ts-comment */
    const filter = activeFilter;
    filter[0][0].autocomplete.push({ name: 'helloservice', value: 'helloservice' });
    component._seqFilters = filter;

    // when
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.mapServiceFilters(metadata);

    // then
    expect(component._seqFilters).toEqual(filter);

    // when
    metadata.filter.services.pop();
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.mapServiceFilters(metadata);

    // then
    expect(component._seqFilters).toEqual(activeFilter);
  });
});
