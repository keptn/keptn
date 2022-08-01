import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbStageOverviewComponent } from './ktb-stage-overview.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { KtbEnvironmentViewModule } from '../ktb-environment-view.module';
import { Project } from '../../../_models/project';
import { DtAutoComplete, DtFilter } from '../../../_models/dt-filter';
import { ApiService } from '../../../_services/api.service';
import { ProjectsMock } from '../../../_services/_mockData/projects.mock';
import { DtFilterFieldChangeEvent } from '@dynatrace/barista-components/filter-field';
import { Stage } from '../../../_models/stage';
import { Service } from '../../../_models/service';

describe('KtbStageOverviewComponent', () => {
  let component: KtbStageOverviewComponent;
  let fixture: ComponentFixture<KtbStageOverviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbEnvironmentViewModule, HttpClientTestingModule, RouterTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbStageOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set service filter on project change', () => {
    // given
    const emitServiceSpy = jest.spyOn(component.filteredServicesChange, 'emit');
    jest
      .spyOn(TestBed.inject(ApiService), 'environmentFilter', 'get')
      .mockReturnValue({ sockshop: { services: ['carts'] } });

    // when
    component.project = Project.fromJSON(ProjectsMock[0]);

    // then
    expect(component._dataSource.data).toEqual({
      autocomplete: [
        {
          name: 'Services',
          autocomplete: [
            {
              name: 'carts-db',
            },
            {
              name: 'carts',
            },
          ],
        } as DtAutoComplete,
      ],
    });
    expect(emitServiceSpy).toHaveBeenCalledWith(['carts']);
    expect(component.filter).toEqual([
      [
        {
          autocomplete: [
            {
              name: 'carts-db',
            },
            {
              name: 'carts',
            },
          ],
          name: 'Services',
        },
        {
          name: 'carts',
        },
      ],
    ]);
  });

  it('should not emit service filter if it is not a new project', () => {
    // given
    component.project = Project.fromJSON(ProjectsMock[0]);

    // when
    const emitServiceSpy = jest.spyOn(component.filteredServicesChange, 'emit');
    component.project = Project.fromJSON(ProjectsMock[0]);

    // then
    expect(emitServiceSpy).not.toHaveBeenCalled();
  });

  it('should emit and save filter on filter change', () => {
    // given
    const emitServiceSpy = jest.spyOn(component.filteredServicesChange, 'emit');

    // when
    component.filterChanged('sockshop', {
      filters: [
        [
          {
            autocomplete: [
              {
                name: 'carts-db',
              },
              {
                name: 'carts',
              },
            ],
          },
          { name: 'carts' },
        ],
        [
          {
            autocomplete: [
              {
                name: 'carts-db',
              },
              {
                name: 'carts',
              },
            ],
          },
          { name: 'carts-db' },
        ],
      ],
    } as DtFilterFieldChangeEvent<DtFilter>);

    // then
    expect(TestBed.inject(ApiService).environmentFilter).toEqual({
      sockshop: {
        services: ['carts', 'carts-db'],
      },
    });
    expect(emitServiceSpy).toHaveBeenCalledWith(['carts', 'carts-db']);
  });

  it('should emit the selected stage information', () => {
    // given
    const mouseEvent = new MouseEvent('click');
    const stopPropagationSpy = jest.spyOn(mouseEvent, 'stopPropagation');
    const emitStageInfoSpy = jest.spyOn(component.selectedStageInfoChange, 'emit');

    // then
    component.selectStage(mouseEvent, { stageName: 'dev' } as Stage, 'approval');

    // then
    expect(stopPropagationSpy).toHaveBeenCalled();
    expect(emitStageInfoSpy).toHaveBeenCalledWith({ stage: { stageName: 'dev' }, filterType: 'approval' });
  });

  it('should return filtered services', () => {
    // given
    jest
      .spyOn(TestBed.inject(ApiService), 'environmentFilter', 'get')
      .mockReturnValue({ sockshop: { services: ['carts'] } });
    component.project = Project.fromJSON(ProjectsMock[0]);

    // when
    const services = component.filterServices([{ serviceName: 'carts' }, { serviceName: 'carts-db' }] as Service[]);

    // then
    expect(services).toEqual([{ serviceName: 'carts' }]);
  });

  it('should return unfiltered services', () => {
    // given
    component.project = Project.fromJSON(ProjectsMock[0]);

    // when
    const services = component.filterServices([{ serviceName: 'carts' }, { serviceName: 'carts-db' }] as Service[]);

    // then
    expect(services).toEqual([{ serviceName: 'carts' }, { serviceName: 'carts-db' }]);
  });
});
