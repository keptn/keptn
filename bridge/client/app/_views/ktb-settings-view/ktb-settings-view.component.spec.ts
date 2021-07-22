import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbSettingsViewComponent } from './ktb-settings-view.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from '@angular/common/http/testing';
import {DataService} from '../../_services/data.service';
import {DataServiceMock} from '../../_services/data.service.mock';
import {of} from 'rxjs';
import {ActivatedRoute, Router} from '@angular/router';

describe('KtbSettingsViewComponent - Create', () => {
  let component: KtbSettingsViewComponent;
  let fixture: ComponentFixture<KtbSettingsViewComponent>;
  let dataService: DataService;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbSettingsViewComponent ],
      imports: [ AppModule, HttpClientTestingModule ],
      providers: [
        {provide: DataService, useClass: DataServiceMock},
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({isCreateMode: true}),
            params: of({}),
            queryParams: of({})
          }
        }
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSettingsViewComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get(DataService);
    fixture.detectChanges();
  });

  it('should create settings view component', () => {
    expect(component).toBeTruthy();
  });

  it('should have create mode enabled when routed to route data contains {isCreateMode: true}', () => {
    expect(component.isCreateMode).toBeTrue();
  });

  it('should have a validation error if project name already exists in projects', async () => {
    // given

    // when
    await dataService.loadProjects();
    component.projectNameControl.setValue('sockshop');

    // then
    expect(component.projectNameControl.hasError('projectName')).toBeTrue();
  });

  it('should navigate to created project', async () => {
    // given
    component.isCreateMode = true;
    component.projectName = 'sockshop';

    // when
    const router = TestBed.inject(Router);
    const routeSpy = spyOn(router, 'navigate');
    await dataService.loadProjects();

    // then
    expect(routeSpy).toHaveBeenCalled();
  });
});

describe('KtbSettingsViewComponent - Edit', () => {
  let component: KtbSettingsViewComponent;
  let fixture: ComponentFixture<KtbSettingsViewComponent>;
  let dataService: DataService;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [KtbSettingsViewComponent],
      imports: [AppModule, HttpClientTestingModule],
      providers: [
        {provide: DataService, useClass: DataServiceMock},
        {
          provide: ActivatedRoute,
          useValue: {
            params: of({projectName: 'sockshop'}),
            data: of({isCreateMode: false}),
            queryParams: of({})
          }
        }
      ]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSettingsViewComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get(DataService);
    fixture.detectChanges();
  });

  it('should create settings view component', () => {
    expect(component).toBeTruthy();
  });

  it('should call DataService.setGitUpstreamUrl on setGitUpstream', () => {
    // given
    const gitData = {
      remoteURI: 'https://test.git',
      gitUser: 'username',
      gitToken: 'token'
    };
    component.projectName = 'sockshop';

    // when
    const spy = spyOn(dataService, 'setGitUpstreamUrl').and.callThrough();
    component.updateGitData(gitData);
    component.setGitUpstream();

    // then
    expect(spy).toHaveBeenCalled();
  });

  it('should have create mode disabled when routed to route data contains {isCreateMode: false}', () => {
    expect(component.isCreateMode).toBeFalse();
  });

  it('should set project name to projectName retrieved by route', () => {
    expect(component.projectName).toEqual('sockshop');
  });

  it('should have an pattern validation error when project name does not match: first letter lowercase, only lowercase, numbers and hyphens allowed', () => {
    // given
    component.isCreateMode = true;

    component.projectNameControl.setValue('Sockshop');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.hasError('pattern'));

    component.projectNameControl.setValue('1ockshop');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.hasError('pattern'));

    component.projectNameControl.setValue('-ockshop');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.hasError('pattern'));

    component.projectNameControl.setValue('$ockshop');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.hasError('pattern'));

    component.projectNameControl.setValue('soCkshop');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.hasError('pattern'));

    component.projectNameControl.setValue('sock_shop');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.hasError('pattern'));

    component.projectNameControl.setValue('sockshop-1');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.errors).toBeNull();
  });

  it('should delete a project and navigate to dashboard', () => {
    // given
    component.projectName = 'sockshop';

    // when
    const router = TestBed.inject(Router);
    const routeSpy = spyOn(router, 'navigate');
    component.deleteProject('sockshop');

    // then
    expect(routeSpy).toHaveBeenCalled();
  });
});


