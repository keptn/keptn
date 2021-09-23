import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceStateInfoComponent } from './ktb-sequence-state-info.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { Project } from '../../_models/project';

describe('KtbSequenceStateInfoComponent', () => {
  let component: KtbSequenceStateInfoComponent;
  let fixture: ComponentFixture<KtbSequenceStateInfoComponent>;
  let dataService: DataService;
  const projectName = 'sockshop';
  let project: Project;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
      providers: [
        {
          provide: DataService,
          useClass: DataServiceMock,
        }
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceStateInfoComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get(DataService);
    dataService.loadProjects(); // reset project.sequences
    // @ts-ignore
    dataService.getProject(projectName).subscribe((pr: Project) => {
      project = pr;
      fixture.detectChanges();
    });
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show sequence info', () => {
    // given
    dataService.loadSequences(project);
    component.sequence = project.sequences[11];
    fixture.detectChanges();

    // then
    const sequenceName = fixture.nativeElement.querySelector('[uitestid=keptn-sequence-info-sequenceName]');
    const serviceName = fixture.nativeElement.querySelector('[uitestid=keptn-sequence-info-serviceName]');
    const status = fixture.nativeElement.querySelector('[uitestid=keptn-sequence-info-status]');

    expect(sequenceName.textContent).toEqual('delivery');
    expect(serviceName.textContent).toEqual('carts:0.12.1');
    expect(status.textContent).toEqual('succeeded');
  });

  it('should show sequence info with 3 stages', () => {
    // given
    dataService.loadSequences(project);
    component.sequence = project.sequences[11];
    fixture.detectChanges();

    // then
    const stageDetails = fixture.nativeElement.querySelectorAll('[uitestid=keptn-sequence-info-stageDetails] ktb-stage-badge');
    expect(stageDetails.length).toEqual(3);

    const firstStage = stageDetails[0].querySelector('dt-tag');
    const secondStage = stageDetails[1].querySelector('dt-tag');
    const thirdStage = stageDetails[2].querySelector('dt-tag');
    expect(firstStage.textContent).toEqual('dev');
    expect(secondStage.textContent).toEqual('staging');
    expect(thirdStage.textContent).toEqual('production');
  });

  it('should show sequence info with 1 stage if showOnlyLastStage is true', () => {
    // given
    dataService.loadSequences(project);
    component.sequence = project.sequences[11];
    component.showOnlyLastStage = true;
    fixture.detectChanges();

    // then
    const stageDetails = fixture.nativeElement.querySelectorAll('[uitestid=keptn-sequence-info-stageDetails] ktb-stage-badge');
    expect(stageDetails.length).toEqual(1);

    const firstStage = stageDetails[0].querySelector('dt-tag');
    expect(firstStage.textContent).toEqual('production');
  });

  it('should trigger click callback on stage', () => {
    // given
    dataService.loadSequences(project);
    const sequence = project.sequences[11]
    component.sequence = sequence;
    const spy = jest.spyOn(component, 'stageClick');
    fixture.detectChanges();

    // when
    const stageDetails = fixture.nativeElement.querySelectorAll('[uitestid=keptn-sequence-info-stageDetails] ktb-stage-badge');
    stageDetails[0].click();

    // then
    expect(spy).toHaveBeenCalledWith(sequence, sequence.getStages()[0]);
  });
});
