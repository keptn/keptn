import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceViewComponent } from './ktb-sequence-view.component';
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
import moment from 'moment';

describe('KtbEventsListComponent', () => {
  let component: KtbSequenceViewComponent;
  let fixture: ComponentFixture<KtbSequenceViewComponent>;

  const projectName = 'sockshop';

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

    // when
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.mapServiceFilters(SequenceMetadataMock);

    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    expect(component.filterFieldData.autocomplete[0].autocomplete).toEqual([
      { name: 'carts-db', value: 'carts-db' },
      { name: 'carts', value: 'carts' },
    ]);
  });

  it('should add a service if it is in a sequence but not in metadata', () => {
    // given
    const metadata = SequenceMetadataMock;
    metadata.filter.services = metadata.filter.services.splice(1, 1); // remove carts-db
    /* eslint-disable @typescript-eslint/ban-ts-comment */
    /* @ts-ignore */ // Ignore private property
    component.project = ProjectsMock[0];
    // @ts-ignore // Ignore private property
    component.project.sequences = SequencesMock;
    // @ts-ignore // Ignore private property
    component.filterFieldData.autocomplete[0].autocomplete = [{ name: 'carts', value: 'carts' }];
    /* eslint-enable @typescript-eslint/ban-ts-comment */

    // when
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.mapServiceFilters(metadata);

    // then
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    expect(component.filterFieldData.autocomplete[0].autocomplete).toEqual([
      { name: 'carts', value: 'carts' },
      { name: 'carts-db', value: 'carts-db' },
    ]);
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
    // @ts-ignore // Ignore private property
    component.filterFieldData.autocomplete[0].autocomplete.push({ name: 'helloservice', value: 'helloservice' });
    /* eslint-enable @typescript-eslint/ban-ts-comment */

    // when
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.mapServiceFilters(metadata);

    // then
    /* eslint-disable @typescript-eslint/ban-ts-comment */
    // @ts-ignore // Ignore private property
    expect(component.filterFieldData.autocomplete[0].autocomplete).toHaveLength(3);

    // As order gets messed up sometimes, it's safer to test each individually

    // @ts-ignore // Ignore private property
    expect(component.filterFieldData.autocomplete[0].autocomplete).toContainEqual({
      name: 'carts-db',
      value: 'carts-db',
    });
    // @ts-ignore // Ignore private property
    expect(component.filterFieldData.autocomplete[0].autocomplete).toContainEqual({ name: 'carts', value: 'carts' });
    // @ts-ignore // Ignore private property
    expect(component.filterFieldData.autocomplete[0].autocomplete).toContainEqual({
      name: 'helloservice',
      value: 'helloservice',
    });
    /* eslint-enable @typescript-eslint/ban-ts-comment */

    // when
    metadata.filter.services.pop();
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    component.mapServiceFilters(metadata);

    // then

    /* eslint-disable @typescript-eslint/ban-ts-comment */
    // @ts-ignore // Ignore private property
    expect(component.filterFieldData.autocomplete[0].autocomplete).toHaveLength(2);
    // @ts-ignore // Ignore private property
    expect(component.filterFieldData.autocomplete[0].autocomplete).toContainEqual({ name: 'carts', value: 'carts' });
    // @ts-ignore // Ignore private property
    expect(component.filterFieldData.autocomplete[0].autocomplete).toContainEqual({
      name: 'carts-db',
      value: 'carts-db',
    });
    /* eslint-enable @typescript-eslint/ban-ts-comment */
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
});
