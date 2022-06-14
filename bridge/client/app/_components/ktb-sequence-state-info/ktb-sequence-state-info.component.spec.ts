import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceStateInfoComponent } from './ktb-sequence-state-info.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { SequencesMock } from '../../_services/_mockData/sequences.mock';

describe('KtbSequenceStateInfoComponent', () => {
  let component: KtbSequenceStateInfoComponent;
  let fixture: ComponentFixture<KtbSequenceStateInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [
        {
          provide: ApiService,
          useClass: ApiServiceMock,
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceStateInfoComponent);
    component = fixture.componentInstance;

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show sequence info', () => {
    // given
    component.sequence = SequencesMock[12];
    fixture.detectChanges();

    // then
    const sequenceName = fixture.nativeElement.querySelector('[uitestid=keptn-sequence-info-sequenceName]');
    const serviceName = fixture.nativeElement.querySelector('[uitestid=keptn-sequence-info-serviceName]');
    const status = fixture.nativeElement.querySelector('[uitestid=keptn-sequence-info-status]');

    expect(sequenceName.textContent).toEqual('delivery');
    expect(serviceName.textContent).toEqual('carts');
    expect(status.textContent).toEqual('succeeded');
  });

  it('should show sequence info with 3 stages', () => {
    // given
    component.sequence = SequencesMock[12];
    fixture.detectChanges();

    // then
    const stageDetails = fixture.nativeElement.querySelectorAll(
      '[uitestid=keptn-sequence-info-stageDetails] ktb-stage-badge'
    );
    expect(stageDetails.length).toEqual(3);

    const firstStage = stageDetails[0].querySelector('dt-tag');
    const secondStage = stageDetails[1].querySelector('dt-tag');
    const thirdStage = stageDetails[2].querySelector('dt-tag');
    expect(firstStage.textContent).toEqual('dev');
    expect(secondStage.textContent).toEqual('staging');
    expect(thirdStage.textContent).toEqual('production');
  });

  it('should show sequence info without stages if showStages is false', () => {
    // given
    component.sequence = SequencesMock[12];
    component.showStages = false;
    fixture.detectChanges();

    // then
    const stageDetails = fixture.nativeElement.querySelectorAll(
      '[uitestid=keptn-sequence-info-stageDetails] ktb-stage-badge'
    );
    expect(stageDetails.length).toEqual(0);
  });

  it('should trigger click callback on stage', () => {
    // given
    const sequence = SequencesMock[12];
    component.sequence = sequence;
    const spy = jest.spyOn(component, 'stageClick');
    fixture.detectChanges();

    // when
    const stageDetails = fixture.nativeElement.querySelectorAll(
      '[uitestid=keptn-sequence-info-stageDetails] ktb-stage-badge'
    );
    stageDetails[0].click();

    // then
    expect(spy).toHaveBeenCalledWith(sequence, sequence?.getStages()[0]);
  });
});
