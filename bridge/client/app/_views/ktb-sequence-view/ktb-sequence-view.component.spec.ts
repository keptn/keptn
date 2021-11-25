import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceViewComponent } from './ktb-sequence-view.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';
import { ProjectsMock } from '../../_services/_mockData/projects.mock';
import { SequencesMock } from '../../_services/_mockData/sequences.mock';
import { POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';

describe('KtbEventsListComponent', () => {
  let httpMock: HttpTestingController;

  let component: KtbSequenceViewComponent;
  let fixture: ComponentFixture<KtbSequenceViewComponent>;
  let dataService: DataService;

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
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceViewComponent);
    component = fixture.componentInstance;

    dataService = fixture.debugElement.injector.get(DataService);
    httpMock = TestBed.inject(HttpTestingController);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show loading indicator while loading', async () => {
    // given
    dataService.loadProjects();
    const projectsRequest = httpMock.expectOne('./api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50');
    projectsRequest.flush({
      projects: ProjectsMock,
      totalCount: ProjectsMock.length,
    });

    fixture.detectChanges();

    const loadingIndicator = fixture.nativeElement.querySelector('[uitestid=keptn-loadingSequences]');
    const emptyStateContainer = fixture.nativeElement.querySelector('[uitestid=keptn-noSequences]');
    const sequenceList = fixture.nativeElement.querySelector('[uitestid=keptn-sequence-view-roots]');

    // then
    expect(loadingIndicator).toBeTruthy();
    expect(emptyStateContainer).toBeFalsy();
    expect(sequenceList).toBeFalsy();
  });

  it('should show empty state if no sequences loaded', async () => {
    // given
    dataService.loadProjects();
    const projectsRequest = httpMock.expectOne('./api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50');
    projectsRequest.flush({
      projects: ProjectsMock,
      totalCount: ProjectsMock.length,
    });

    const sequencesRequest = httpMock.expectOne(`./api/controlPlane/v1/sequence/${projectName}?pageSize=25`);
    sequencesRequest.flush({
      states: [],
      totalCount: 0,
    });

    httpMock.verify();
    fixture.detectChanges();

    const loadingIndicator = fixture.nativeElement.querySelector('[uitestid=keptn-loadingSequences]');
    const emptyStateContainer = fixture.nativeElement.querySelector('[uitestid=keptn-noSequences]');
    const sequenceList = fixture.nativeElement.querySelector('[uitestid=keptn-sequence-view-roots]');

    // then
    expect(loadingIndicator).toBeFalsy();
    expect(emptyStateContainer).toBeTruthy();
    expect(sequenceList).toBeFalsy();
  });

  it('should show list of sequences', async () => {
    // given
    dataService.loadProjects();
    const projectsRequest = httpMock.expectOne('./api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50');
    projectsRequest.flush({
      projects: ProjectsMock,
      totalCount: ProjectsMock.length,
    });

    const sequencesRequest = httpMock.expectOne(`./api/controlPlane/v1/sequence/${projectName}?pageSize=25`);
    sequencesRequest.flush({
      states: SequencesMock,
      totalCount: SequencesMock.length,
    });

    httpMock.verify();
    fixture.detectChanges();

    const loadingIndicator = fixture.nativeElement.querySelector('[uitestid=keptn-loadingSequences]');
    const emptyStateContainer = fixture.nativeElement.querySelector('[uitestid=keptn-noSequences]');
    const sequenceList = fixture.nativeElement.querySelector('[uitestid=keptn-sequence-view-roots]');

    // then
    expect(loadingIndicator).toBeFalsy();
    expect(emptyStateContainer).toBeFalsy();
    expect(sequenceList).toBeTruthy();
  });

  it('should show empty list after filter applied', async () => {
    // given
    dataService.loadProjects();
    const projectsRequest = httpMock.expectOne('./api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50');
    projectsRequest.flush({
      projects: ProjectsMock,
      totalCount: ProjectsMock.length,
    });

    const sequencesRequest = httpMock.expectOne(`./api/controlPlane/v1/sequence/${projectName}?pageSize=25`);
    sequencesRequest.flush({
      states: SequencesMock,
      totalCount: SequencesMock.length,
    });

    httpMock.verify();
    fixture.detectChanges();

    // when
    clickFilterCheckbox('carts-db');
    clickFilterCheckbox('delivery');
    fixture.detectChanges();

    // then
    const loadingIndicator = fixture.nativeElement.querySelector('[uitestid=keptn-loadingSequences]');
    const emptyStateFilteredContainer = fixture.nativeElement.querySelector('[uitestid=keptn-noSequencesFiltered]');
    const sequenceList = fixture.nativeElement.querySelector('[uitestid=keptn-sequence-view-roots]');

    expect(loadingIndicator).toBeFalsy();
    expect(emptyStateFilteredContainer).toBeTruthy();
    expect(sequenceList).toBeFalsy();
  });

  function clickFilterCheckbox(selector: string): void {
    const checkboxes = document.evaluate(
      `//dt-checkbox[contains(., '${selector}')]`,
      document,
      null,
      XPathResult.ANY_TYPE,
      null
    );
    const checkbox = checkboxes.iterateNext();
    if (checkbox) {
      checkbox.childNodes[0].childNodes[0].dispatchEvent(new Event('click'));
    }
  }
});
