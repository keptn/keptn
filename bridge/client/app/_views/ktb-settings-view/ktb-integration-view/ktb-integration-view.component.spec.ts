import { ComponentFixture, TestBed } from '@angular/core/testing';
import { compare, KtbIntegrationViewComponent, sortRegistrations } from './ktb-integration-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, convertToParamMap, ParamMap } from '@angular/router';
import { BehaviorSubject, firstValueFrom } from 'rxjs';
import { ApiService } from '../../../_services/api.service';
import { ApiServiceMock } from '../../../_services/api.service.mock';
import { RouterTestingModule } from '@angular/router/testing';
import { KtbIntegrationViewModule } from './ktb-integration-view.module';
import { UniformRegistration } from '../../../_models/uniform-registration';
import { UniformSubscription } from '../../../_models/uniform-subscription';
import { ElementRef, EmbeddedViewRef, TemplateRef } from '@angular/core';

class MockElementRef extends ElementRef {
  nativeElement = {};

  constructor() {
    super(null);
  }
}

class MockTemplateRef extends TemplateRef<unknown> {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  createEmbeddedView(_context: unknown): EmbeddedViewRef<unknown> {
    return {} as unknown as EmbeddedViewRef<unknown>;
  }

  get elementRef(): ElementRef {
    return new MockElementRef();
  }
}

describe(KtbIntegrationViewComponent.name, () => {
  const projectName = 'sockshop';
  let component: KtbIntegrationViewComponent;
  let fixture: ComponentFixture<KtbIntegrationViewComponent>;
  let paramsSubject: BehaviorSubject<ParamMap>;

  beforeEach(async () => {
    paramsSubject = new BehaviorSubject(convertToParamMap({}));
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbIntegrationViewModule, HttpClientTestingModule, RouterTestingModule],
      providers: [
        { provide: ApiService, useClass: ApiServiceMock },
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: paramsSubject.asObservable(),
          },
        },
      ],
    }).compileComponents();

    localStorage.setItem('keptn_integration_dates', '');
    fixture = TestBed.createComponent(KtbIntegrationViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should provide the project path param', async () => {
    // given
    paramsSubject.next(convertToParamMap({ projectName }));

    // when
    const actual = await firstValueFrom(component.params$);

    // then
    expect(actual.projectName).toBe('sockshop');
  });

  it('should provide the project path param and the integration param', async () => {
    // given
    paramsSubject.next(
      convertToParamMap({
        projectName,
        integrationId: 'abcdef',
      })
    );

    // when
    const actual = await firstValueFrom(component.params$);

    // then
    expect(actual.projectName).toBe('sockshop');
    expect(actual.integrationId).toBe('abcdef');
  });

  it('should emit the selected uniform registration when the path param is provided', async () => {
    // given
    paramsSubject.next(
      convertToParamMap({
        projectName,
        integrationId: 'abcdef',
      })
    );

    // when
    const actual = await firstValueFrom(component.selectedUniformRegistrationId$);

    // then
    expect(actual).toBe('abcdef');
  });

  it('should load uniform registrations and build the data source', async () => {
    // given
    // when
    const dataSource = await firstValueFrom(component.uniformRegistrations$);

    // then
    expect(dataSource.data).toBeTruthy();
    expect(dataSource.data.length).toBe(11);
  });

  it('should select a uniform registration by the user', (done) => {
    // given
    const id = 'af567de129';
    const u = new UniformRegistration();
    u.id = id;
    component.selectedUniformRegistrationId$.subscribe((actual) => {
      // then
      expect(actual).toBe(id);
      done();
    });

    // when
    component.setSelectedUniformRegistration(u, projectName);
  });

  it('should pre select a uniform registration when the path param is set', async () => {
    // given
    const integrationId = 'keptn-lighthouse-service-8feec7146c19fa08bd65664b8d47f153';
    paramsSubject.next(
      convertToParamMap({
        projectName,
        integrationId,
      })
    );

    // when
    const actual = await firstValueFrom(component.selectedUniformRegistration$);

    // then
    expect(actual).toBeTruthy();
    expect(actual?.id).toBe(integrationId);
  });

  it('should pre load logs', async () => {
    // given
    const integrationId = 'keptn-lighthouse-service-8feec7146c19fa08bd65664b8d47f153';
    paramsSubject.next(
      convertToParamMap({
        projectName,
        integrationId,
      })
    );

    // when
    const actual = await firstValueFrom(component.uniformRegistrationLogs$);

    // then
    expect(actual).toBeTruthy();
    expect(actual.length).toBe(10);
  });

  it('should get subscriptions for a uniform registration', () => {
    // given
    const u = new UniformRegistration();
    const s1 = new UniformSubscription();
    s1.event = 'Dog';
    const s2 = new UniformSubscription();
    s2.event = 'Cat';
    u.subscriptions = [s1, s2];

    // when
    const actual = component.getSubscriptions(u, projectName);

    // then
    expect(actual).toEqual(['Dog', 'Cat']);
  });

  it('should convert unknown to a uniform registration', () => {
    // given
    const u: unknown = new UniformRegistration();

    // when
    const actual = component.toUniformRegistration(u);

    // then
    expect(actual).toBeTruthy();
  });

  it('should return an overlay', () => {
    // given
    const templateRef = new MockTemplateRef();

    // when
    const actual = component.getOverlay(new UniformRegistration(), projectName, templateRef);

    // then
    expect(actual).toEqual(templateRef);
  });

  describe(KtbIntegrationViewComponent.name + 'HelperFunctions', () => {
    it('should compare', () => {
      expect(compare('a', 'b', true)).toBe(-1);
      expect(compare('a', 'b', false)).toBe(1);
      expect(compare('b', 'a', true)).toBe(1);
      expect(compare('b', 'a', false)).toBe(-1);
      expect(compare('a', 'a', true)).toBe(0);
      expect(compare('a', 'a', false)).toBe(-0);
    });

    it('should sort registrations by name (default)', () => {
      // given
      const uniforms = getUniforms();

      // when
      const actualAsc = sortRegistrations(uniforms, 'anythingDoesNotMatter', true);
      const actualDesc = sortRegistrations(uniforms, 'anythingDoesNotMatter', false);

      // then
      const idSortOrder = (r: UniformRegistration): string => r.id;
      expect(actualAsc.map(idSortOrder)).toEqual(['2', '1', '3']);
      expect(actualDesc.map(idSortOrder)).toEqual(['3', '1', '2']);
    });

    it('should sort registrations by host', () => {
      // given
      const uniforms = getUniforms();

      // when
      const actualAsc = sortRegistrations(uniforms, 'host', true);
      const actualDesc = sortRegistrations(uniforms, 'host', false);

      // then
      const idSortOrder = (r: UniformRegistration): string => r.id;
      expect(actualAsc.map(idSortOrder)).toEqual(['3', '2', '1']);
      expect(actualDesc.map(idSortOrder)).toEqual(['1', '2', '3']);
    });

    it('should sort registrations by location', () => {
      // given
      const uniforms = getUniforms();

      // when
      const actualAsc = sortRegistrations(uniforms, 'location', true);
      const actualDesc = sortRegistrations(uniforms, 'location', false);

      // then
      const idSortOrder = (r: UniformRegistration): string => r.id;
      expect(actualAsc.map(idSortOrder)).toEqual(['1', '2', '3']);
      expect(actualDesc.map(idSortOrder)).toEqual(['3', '2', '1']);
    });

    it('should sort registrations by namespace', () => {
      // given
      const uniforms = getUniforms();

      // when
      const actualAsc = sortRegistrations(uniforms, 'namespace', true);
      const actualDesc = sortRegistrations(uniforms, 'namespace', false);

      // then
      const idSortOrder = (r: UniformRegistration): string => r.id;
      expect(actualAsc.map(idSortOrder)).toEqual(['2', '3', '1']);
      expect(actualDesc.map(idSortOrder)).toEqual(['1', '3', '2']);
    });

    function getUniform(id: string, name: string, host: string, loc: string, namespace: string): UniformRegistration {
      const u1 = new UniformRegistration();
      u1.id = id;
      u1.name = name;
      u1.metadata = {
        hostname: host,
        location: loc,
        deplyomentname: '',
        distributorversion: '',
        integrationversion: '',
        status: '',
        lastseen: '',
        kubernetesmetadata: {
          namespace,
          podname: '',
          deploymentname: '',
        },
      };
      return u1;
    }

    function getUniforms(): UniformRegistration[] {
      return [
        getUniform('1', 'B', 'C', 'A', 'C'),
        getUniform('2', 'A', 'B', 'B', 'A'),
        getUniform('3', 'C', 'A', 'C', 'B'),
      ];
    }
  });
});
