import { ComponentFixture, fakeAsync, TestBed, tick } from '@angular/core/testing';
import { KtbCreateSecretFormComponent } from './ktb-create-secret-form.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { DataService } from '../../_services/data.service';
import { Secret } from '../../_models/secret';
import { of, throwError } from 'rxjs';
import { Router } from '@angular/router';
import { SecretScope } from '../../../../shared/interfaces/secret-scope';
import { POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';

describe('KtbCreateSecretFormComponent with valid scopes', () => {
  let component: KtbCreateSecretFormComponent;
  let fixture: ComponentFixture<KtbCreateSecretFormComponent>;
  let dataService: DataService;
  let router: Router;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [{ provide: POLLING_INTERVAL_MILLIS, useValue: 0 }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbCreateSecretFormComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get(DataService);
    router = TestBed.inject(Router);
    httpMock = TestBed.inject(HttpTestingController);
    httpMock.expectOne('./api/secrets/v1/scope').flush({
      scopes: [SecretScope.WEBHOOK, SecretScope.DEFAULT, 'dynatrace-service'],
    });
  });

  it('should create', () => {
    expect(component).toBeTruthy();
    component.ngOnDestroy();
  });

  it('should create secret with keptn-default scope', () => {
    // given
    const spy = jest.spyOn(dataService, 'addSecret');
    const secret = insertDefaultSecret(component);
    component.createSecretForm.updateValueAndValidity();
    fixture.detectChanges();

    expect(component.createSecretForm.errors).toBeNull();
    expect(component.isFormValid()).toBe(true);
    component.createSecret();

    // then
    expect(spy).toHaveBeenCalledWith(secret);
    component.ngOnDestroy();
  });

  it('should create secret with selected scope', () => {
    // given
    const spy = jest.spyOn(dataService, 'addSecret').mockReturnValue(of({}));
    const routerSpy = jest.spyOn(router, 'navigate');
    const secret: Secret = new Secret();
    secret.name = 'test';
    secret.scope = component.scopes?.[1] ?? '';
    secret.data?.push({ key: 'testKey', value: 'testValue' });

    // when
    component.nameControl.setValue(secret.name);
    component.scopeControl.setValue(secret.scope);
    if (secret.data) {
      component.dataControl.controls[0].get('key')?.setValue(secret.data[0].key);
      component.dataControl.controls[0].get('value')?.setValue(secret.data[0].value);
    }
    component.createSecretForm.updateValueAndValidity();
    fixture.detectChanges();

    expect(component.createSecretForm.errors).toBeNull();
    expect(component.isFormValid()).toBe(true);
    component.createSecret();

    // then
    expect(spy).toHaveBeenCalledWith(secret);
    expect(routerSpy).toHaveBeenCalled();
    expect(component.isLoading).toBe(false);
    component.ngOnDestroy();
  });

  it('should handle failed creating secret', () => {
    // given
    const spy = jest.spyOn(dataService, 'addSecret').mockReturnValue(throwError({}));
    const routerSpy = jest.spyOn(router, 'navigate');
    const secret: Secret = new Secret();
    secret.name = 'test';
    secret.scope = component.scopes?.[1] ?? '';
    secret.data?.push({ key: 'testKey', value: 'testValue' });

    // when
    component.nameControl.setValue(secret.name);
    component.scopeControl.setValue(secret.scope);
    if (secret.data) {
      component.dataControl.controls[0].get('key')?.setValue(secret.data[0].key);
      component.dataControl.controls[0].get('value')?.setValue(secret.data[0].value);
    }
    component.createSecretForm.updateValueAndValidity();
    fixture.detectChanges();

    expect(component.createSecretForm.errors).toBeNull();
    expect(component.isFormValid()).toBe(true);
    component.createSecret();

    // then
    expect(spy).toHaveBeenCalledWith(secret);
    expect(routerSpy).not.toHaveBeenCalled();
    expect(component.isLoading).toBe(false);
    component.ngOnDestroy();
  });

  it('should add key/value pair', () => {
    // when
    component.addPair();

    // then
    expect(component.dataControl.controls.length).toBe(2);
    component.ngOnDestroy();
  });

  it('should remove key/value pair', () => {
    // when
    component.addPair();
    component.removePair(0);

    // then
    expect(component.dataControl.controls.length).toBe(1);
    component.ngOnDestroy();
  });
});
describe('KtbCreateSecretFormComponent scopes', () => {
  let component: KtbCreateSecretFormComponent;
  let fixture: ComponentFixture<KtbCreateSecretFormComponent>;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();
    httpMock = TestBed.inject(HttpTestingController);
  });

  it('should have scopes set to undefined and running retry interval', fakeAsync(() => {
    createComponent(); // has to be inside fakeAsync because the interval is created in the constructor
    const loadingSpy = jest.spyOn(component, 'isLoading', 'set');
    tick();
    httpMock.expectOne('./api/secrets/v1/scope').error(new ErrorEvent('error'));
    expect(component.isLoading).toBe(false);
    expect(loadingSpy).toHaveBeenCalledWith(true);
    expect(loadingSpy).toHaveBeenCalledWith(false);
    expect(component.scopes).toBe(undefined);
    component.ngOnDestroy();
  }));

  it('should have invalid form, if scopes are loading', fakeAsync(() => {
    createComponent();
    tick();
    expect(component.isLoading).toBe(true);
    insertDefaultSecret(component);
    expect(component.isFormValid()).toBe(false);
    component.ngOnDestroy();
  }));

  it('should have scopes set to array', fakeAsync(() => {
    createComponent();
    const scopes = [SecretScope.WEBHOOK];
    tick();
    httpMock.expectOne('./api/secrets/v1/scope').flush({ scopes }, { status: 200, statusText: 'OK' });
    fixture.detectChanges();
    expect(component.scopes).toEqual(scopes);
    component.ngOnDestroy();
  }));

  function createComponent(): void {
    fixture = TestBed.createComponent(KtbCreateSecretFormComponent);
    component = fixture.componentInstance;
  }
});

function insertDefaultSecret(component: KtbCreateSecretFormComponent): Secret {
  const secret: Secret = new Secret();
  secret.name = 'test';
  secret.scope = SecretScope.DEFAULT;
  secret.data?.push({ key: 'testKey', value: 'testValue' });

  // when
  component.nameControl.setValue(secret.name);
  component.scopeControl.setValue(secret.scope);
  if (secret.data) {
    component.dataControl.controls[0].get('key')?.setValue(secret.data[0].key);
    component.dataControl.controls[0].get('value')?.setValue(secret.data[0].value);
  }
  return secret;
}
