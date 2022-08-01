import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbStageDetailsComponent } from './ktb-stage-details.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbEnvironmentViewModule } from '../ktb-environment-view.module';
import { RouterTestingModule } from '@angular/router/testing';
import { Stage } from '../../../_models/stage';
import { DtToggleButtonChange } from '@dynatrace/barista-components/toggle-button-group';
import { Service } from '../../../_models/service';

describe('KtbStageDetailsComponent', () => {
  let component: KtbStageDetailsComponent;
  let fixture: ComponentFixture<KtbStageDetailsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbEnvironmentViewModule, HttpClientTestingModule, RouterTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbStageDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set filterType and stageInfo', () => {
    // given
    const stageInfoChangeSpy = jest.spyOn(component.selectedStageInfoChange, 'emit');

    // when
    component.selectedStageInfo = { stage: { stageName: 'dev' } as Stage, filterType: 'approval' };

    // then
    expect(component.filterEventType).toBe('approval');
    expect(component.selectedStageInfo).toEqual({ stage: { stageName: 'dev' } as Stage, filterType: 'approval' });
    expect(stageInfoChangeSpy).not.toHaveBeenCalled();
  });

  it('should set and emit filter event type', () => {
    // given
    const stageInfoChangeSpy = jest.spyOn(component.selectedStageInfoChange, 'emit');

    // when
    component.selectFilterEvent(
      { stageName: 'dev' } as Stage,
      {
        isUserInput: true,
        value: 'approval',
        source: {
          get selected(): boolean {
            return true;
          },
        },
      } as DtToggleButtonChange<unknown>
    );

    // then
    expect(stageInfoChangeSpy).toHaveBeenCalledWith({ stage: { stageName: 'dev' }, filterType: 'approval' });
  });

  it('should set and emit empty filter event type', () => {
    // given
    const stageInfoChangeSpy = jest.spyOn(component.selectedStageInfoChange, 'emit');

    // when
    component.selectFilterEvent(
      { stageName: 'dev' } as Stage,
      {
        isUserInput: true,
        value: 'approval',
        source: {
          get selected(): boolean {
            return false;
          },
        },
      } as DtToggleButtonChange<unknown>
    );

    // then
    expect(stageInfoChangeSpy).toHaveBeenCalledWith({ stage: { stageName: 'dev' }, filterType: undefined });
  });

  it('should return service link', () => {
    // given, when
    const link = component.getServiceLink(
      {
        serviceName: 'carts',
        stage: 'dev',
        get deploymentContext(): string | undefined {
          return 'keptnContext';
        },
      } as Service,
      'sockshop'
    );

    expect(link).toEqual(['/project', 'sockshop', 'service', 'carts', 'context', 'keptnContext', 'stage', 'dev']);
  });

  it('should filter services for certain filterType and should not reset filter', () => {
    // given
    component.filteredServices = ['carts'];
    component.filterEventType = 'approval';

    // when
    const services = component.filterServices(
      { stageName: 'dev' } as Stage,
      [
        {
          serviceName: 'carts',
        },
        {
          serviceName: 'carts-db',
        },
      ] as Service[],
      'approval'
    );

    // then
    expect(services).toEqual([{ serviceName: 'carts' }]);
    expect(component.filterEventType).toBe('approval');
  });

  it('should reset filter if the type does not have any services', () => {
    const stageInfoChangeSpy = jest.spyOn(component.selectedStageInfoChange, 'emit');
    component.filteredServices = ['carts']; // => carts does not have an approval
    component.filterEventType = 'approval';

    // when
    const services = component.filterServices(
      { stageName: 'dev' } as Stage,
      [
        {
          serviceName: 'carts-db',
        },
      ] as Service[],
      'approval'
    );

    // then
    expect(services).toEqual([]);
    expect(component.filterEventType).toBeUndefined();
    expect(stageInfoChangeSpy).toHaveBeenCalledWith({ stage: { stageName: 'dev' }, filterType: undefined });
    expect(component.selectedStageInfo).toEqual({ stage: { stageName: 'dev' }, filterType: undefined });
  });
});
