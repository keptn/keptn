import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbCreateSecretFormComponent } from './ktb-create-secret-form.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataService } from '../../_services/data.service';

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

    // when
    component.getFormControl('name')?.setValue('test');
    component.data?.controls[0].get('key')?.setValue('testKey');
    component.data?.controls[0].get('value')?.setValue('testValue');
    createButton.click();

    // then
    expect(spy).toHaveBeenCalled();
    const argument = spy.mock.calls[0][0];
    expect(argument.name).toEqual('test');
    expect(argument.scope).toEqual(component.scopes[0]);
    expect(argument.data[0].key).toEqual('testKey');
    expect(argument.data[0].value).toEqual('testValue');
  });

  it('should create secret with selected scope', () => {
    // given
    const createButton = fixture.nativeElement.querySelector('[uitestid=keptn-secret-create-button]');
    const spy = jest.spyOn(dataService, 'addSecret');

    // when
    component.getFormControl('name')?.setValue('test');
    component.getFormControl('scope')?.setValue(component.scopes[1]);
    component.data?.controls[0].get('key')?.setValue('testKey');
    component.data?.controls[0].get('value')?.setValue('testValue');
    createButton.click();

    // then
    expect(spy).toHaveBeenCalled();
    const argument = spy.mock.calls[0][0];
    expect(argument.name).toEqual('test');
    expect(argument.scope).toEqual(component.scopes[1]);
    expect(argument.data[0].key).toEqual('testKey');
    expect(argument.data[0].value).toEqual('testValue');
  });
});
