import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbCreateSecretFormComponent } from './ktb-create-secret-form.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataService } from '../../_services/data.service';
import { Secret } from '../../_models/secret';

describe('KtbCreateSecretFormComponent', () => {
  let component: KtbCreateSecretFormComponent;
  let fixture: ComponentFixture<KtbCreateSecretFormComponent>;
  let dataService: DataService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbCreateSecretFormComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get(DataService);
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should create secret with keptn-default scope', () => {
    // given
    const createButton = fixture.nativeElement.querySelector('[uitestid=keptn-secret-create-button]');
    const spy = jest.spyOn(dataService, 'addSecret');
    const secret: Secret = new Secret();
    secret.name = 'test';
    secret.scope = component.scopes[0];
    secret.data.push({key: 'testKey', value: 'testValue'});

    // when
    component.getFormControl('name')?.setValue(secret.name);
    component.data?.controls[0].get('key')?.setValue(secret.data[0].key);
    component.data?.controls[0].get('value')?.setValue(secret.data[0].value);
    createButton.click();

    // then
    expect(spy).toHaveBeenCalledWith(secret);
  });

  it('should create secret with selected scope', () => {
    // given
    const createButton = fixture.nativeElement.querySelector('[uitestid=keptn-secret-create-button]');
    const spy = jest.spyOn(dataService, 'addSecret');
    const secret: Secret = new Secret();
    secret.name = 'test';
    secret.scope = component.scopes[1];
    secret.data.push({key: 'testKey', value: 'testValue'});

    // when
    component.getFormControl('name')?.setValue(secret.name);
    component.getFormControl('scope')?.setValue(secret.scope);
    component.data?.controls[0].get('key')?.setValue(secret.data[0].key);
    component.data?.controls[0].get('value')?.setValue(secret.data[0].value);
    createButton.click();

    // then
    expect(spy).toHaveBeenCalledWith(secret);
  });

  it('remove key/value pair should be disabled', () => {
    // given
    fixture.detectChanges();
    const removePairButton = fixture.nativeElement.querySelector('[uitestid=keptn-secret-remove-pair-button]');

    // then
    expect(removePairButton.disabled).toBe(true);
  });

  it('should add key/value pair', () => {
    // given
    const addPairButton = fixture.nativeElement.querySelector('[uitestid=keptn-secret-add-pair-button]');

    // when
    addPairButton.click();

    // then
    expect(component.dataControls.length).toBe(2);
  });

  it('should remove key/value pair', () => {
    // given
    const addPairButton = fixture.nativeElement.querySelector('[uitestid=keptn-secret-add-pair-button]');

    // when
    addPairButton.click();
    fixture.detectChanges();
    const removePairButtons: HTMLElement[] = Array.from(fixture.nativeElement.querySelectorAll('[uitestid=keptn-secret-remove-pair-button]'));
    removePairButtons[0].click();

    // then
    expect(component.dataControls.length).toBe(1);
  });
});
