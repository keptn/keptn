import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbWebhookSettingsComponent } from './ktb-webhook-settings.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AbstractControl } from '@angular/forms';
import { WebhookConfigMock } from '../../_services/_mockData/webhook-config.mock';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { UniformRegistrationsMock } from '../../_models/uniform-registrations.mock';
import { DataService } from '../../_services/data.service';
import { throwError } from 'rxjs';
import { HttpErrorResponse } from '@angular/common/http';
import { DataServiceMock } from '../../_services/data.service.mock';

describe('KtbWebhookSettingsComponent', () => {
  let component: KtbWebhookSettingsComponent;
  let fixture: ComponentFixture<KtbWebhookSettingsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
      providers: [
        {provide: DataService, useClass: DataServiceMock},
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbWebhookSettingsComponent);
    component = fixture.componentInstance;
    component.subscription = new UniformSubscription();
    component.subscriptionExists = false;
    component.projectName = 'sockshop';
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should be invalid URL', () => {
    // @ts-ignore
    const urlControl: AbstractControl = component.webhookConfigForm.get('url');
    const urls = ['', '://keptn.sh', 'keptnsh', 'keptn@sh.sh', 'keptn:sh'];

    for (const url of urls) {
      urlControl.setValue(url);
      expect(urlControl.valid).toEqual(false);
    }
  });

  it('should be valid URL', () => {
    // @ts-ignore
    const urlControl: AbstractControl = component.webhookConfigForm.get('url');
    const urls = ['https://keptn.sh', 'http://keptn.sh', 'http://www.keptn.sh', 'keptn.sh', 'keptn.sh/#id', 'keptn.sh/sh/', 'www.keptn.sh'];

    for (const url of urls) {
      urlControl.setValue(url);
      expect(urlControl.valid).toEqual(true);
    }
  });

  it('should add headers', () => {
    // given
    const addHeaderButton = getAddHeaderButton();
    let headers = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form');
    expect(headers.length).toEqual(0);

    // then
    for (let i = 1; i <= 2; ++i) {
      addHeaderButton.click();
      fixture.detectChanges();
      headers = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form');
      expect(headers.length).toEqual(i);
    }
  });

  it('should remove headers', () => {
    // given
    setParameters();
    fixture.detectChanges();

    // when
    let removeButtons = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form button');
    const lengthBefore = removeButtons.length;
    removeButtons[0].click();
    fixture.detectChanges();

    // then
    removeButtons = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form button');
    expect(lengthBefore - 1).toEqual(removeButtons.length);
  });

  it('should fill data', () => {
    // given
    setParameters();
    fixture.detectChanges();

    // then
    const headers = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form');
    expect(component.getFormControl('url').value).toEqual(WebhookConfigMock.url);
    expect(component.getFormControl('method').value).toEqual(WebhookConfigMock.method);
    expect(component.getFormControl('payload').value).toEqual(WebhookConfigMock.payload);
    expect(component.getFormControl('proxy').value).toEqual(WebhookConfigMock.proxy);

    for (let i = 0; i < component.headerControls.length; ++i) {
      expect(WebhookConfigMock.header[i]).toEqual({
        name: component.headerControls[i].get('name')?.value,
        value: component.headerControls[i].get('value')?.value,
      });
    }
    expect(headers.length).toEqual(WebhookConfigMock.header.length);
  });

  it('should revert loading on fail', () => {
    const dataService = TestBed.inject(DataService);
    const getWebhookSpy = jest.spyOn(dataService, 'getWebhookConfig');
    getWebhookSpy.mockReturnValue(throwError(new HttpErrorResponse({error: ''})));
    setParameters();
    fixture.detectChanges();
    expect(getWebhookSpy).toHaveBeenCalled();
    expect(component.loading).toEqual(false);
  });

  function setParameters(): void {
    component.subscription = UniformRegistrationsMock[2].subscriptions[0];
    component.subscriptionExists = true;
    component.projectName = 'sockshop';
  }

  it('should be invalid form', () => {
    // given
    const urlControl = component.getFormControl('url');
    const methodControl = component.getFormControl('method');
    const payload = component.getFormControl('payload');
    const proxyControl = component.getFormControl('proxy');

    expect(component.webhookConfigForm.valid).toEqual(false);

    // when
    urlControl.setValue('keptn.sh');
    // then
    expect(component.webhookConfigForm.valid).toEqual(false);

    // when
    methodControl.setValue('POST');
    expect(component.webhookConfigForm.valid).toEqual(true);

    // when
    payload.setValue('{}');
    // then
    expect(component.webhookConfigForm.valid).toEqual(true);

    // when
    proxyControl.setValue('keptn');
    // then
    expect(component.webhookConfigForm.valid).toEqual(false);

    // when
    proxyControl.setValue('keptn.sh');
    const addHeaderButton = getAddHeaderButton();
    addHeaderButton.click();
    component.addHeader();

    // then
    expect(component.webhookConfigForm.valid).toEqual(false);
  });

  it('should be valid form', () => {
    // given
    const urlControl = component.getFormControl('url');
    const methodControl = component.getFormControl('method');
    const payload = component.getFormControl('payload');
    const proxyControl = component.getFormControl('proxy');

    // when
    urlControl.setValue('https://keptn.sh');
    methodControl.setValue('POST');
    payload.setValue('{}');

    // then
    expect(component.webhookConfigForm.valid).toEqual(true);

    // when
    proxyControl.setValue('https://keptn.sh');

    // then
    expect(component.webhookConfigForm.valid).toEqual(true);

    // when
    component.addHeader('content-type', 'application/json');

    // then
    expect(component.webhookConfigForm.valid).toEqual(true);
  });

  function getAddHeaderButton(): HTMLElement {
    return fixture.nativeElement.querySelector('[uitestid="ktb-webhook-settings-add-header-button"]');
  }
});
