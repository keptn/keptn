import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbWebhookSettingsComponent } from './ktb-webhook-settings.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AbstractControl } from '@angular/forms';
import { WebhookConfigMock } from '../../_services/_mockData/webhook-config.mock';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';

describe('KtbWebhookSettingsComponent', () => {
  let component: KtbWebhookSettingsComponent;
  let fixture: ComponentFixture<KtbWebhookSettingsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
      providers: [{ provide: DataService, useClass: DataServiceMock }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbWebhookSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should be invalid URL when not starting with http(s)://', () => {
    const urlControl: AbstractControl = component.webhookConfigForm.get('url') as AbstractControl;
    const urls = [
      '://keptn.sh',
      'www.keptn.sh',
      'http:/www.keptn.sh',
      'http//www.keptn.sh',
      'htp://www.keptn.sh',
      'ftp://www.keptn.sh',
    ];

    for (const url of urls) {
      urlControl.setValue(url);
      expect(urlControl.valid).toEqual(false);
    }
  });

  it('should be valid URL when it starts with http(s)://', () => {
    const urlControl: AbstractControl = component.webhookConfigForm.get('url') as AbstractControl;
    const urls = [
      'https://keptn.sh',
      'http://keptn.sh',
      'http://www.keptn.sh',
      'http://my-jenkins-servcer.default.svc.cluster.local:8080/job/nodejs example app/build?token=',
      'http://my-jenkins-servcer.default.svc.cluster.local:8080/job/nodejs%20example%20app/build?token=',
      'http://my-jenkins-servcer.com/job/nodejs example app/build?token=',
    ];

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
    let headerRows = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form');
    const lengthBefore = headerRows.length;
    const buttons = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form button');

    buttons[1].click();
    fixture.detectChanges();

    // then
    headerRows = fixture.nativeElement.querySelectorAll('div[formarrayname="header"] form');
    expect(lengthBefore - 1).toEqual(headerRows.length);
  });

  it('should fill form fields with provided data', () => {
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

  function setParameters(): void {
    component.webhook = WebhookConfigMock;
  }

  it('should be invalid form if only URL is set', () => {
    // given
    const urlControl = component.getFormControl('url');
    urlControl.setValue('keptn.sh');
    expect(component.webhookConfigForm.valid).toEqual(false);
  });

  it('should be invalid form if proxy is invalid', () => {
    // given
    const urlControl = component.getFormControl('url');
    const methodControl = component.getFormControl('method');
    const proxyControl = component.getFormControl('proxy');

    // when
    urlControl.setValue('keptn.sh');
    methodControl.setValue('POST');
    proxyControl.setValue('keptn');
    // then
    expect(component.webhookConfigForm.valid).toEqual(false);
  });

  it('should be invalid form if it has empty header configuration', () => {
    // given
    const urlControl = component.getFormControl('url');
    const methodControl = component.getFormControl('method');

    // when
    urlControl.setValue('keptn.sh');
    methodControl.setValue('POST');

    const addHeaderButton = getAddHeaderButton();
    addHeaderButton.click();
    component.addHeader();

    // then
    expect(component.webhookConfigForm.valid).toEqual(false);
  });

  it('should be valid form if URL and method is set', () => {
    // given
    const urlControl = component.getFormControl('url');
    const methodControl = component.getFormControl('method');

    // when
    urlControl.setValue('https://keptn.sh');
    methodControl.setValue('POST');
    expect(component.webhookConfigForm.valid).toEqual(true);
  });

  it('should be valid form if valid proxy is set', () => {
    // given
    const urlControl = component.getFormControl('url');
    const methodControl = component.getFormControl('method');
    const proxyControl = component.getFormControl('proxy');

    // when
    urlControl.setValue('https://keptn.sh');
    methodControl.setValue('POST');
    proxyControl.setValue('https://keptn.sh');

    // then
    expect(component.webhookConfigForm.valid).toEqual(true);
  });

  it('should be valid form if valid header is added', () => {
    // given
    const urlControl = component.getFormControl('url');
    const methodControl = component.getFormControl('method');

    // when
    urlControl.setValue('https://keptn.sh');
    methodControl.setValue('POST');

    // when
    component.addHeader('content-type', 'application/json');

    // then
    expect(component.webhookConfigForm.valid).toEqual(true);
  });

  it('should emit the set values in the form', () => {
    // given
    const spy = jest.spyOn(component.webhookChange, 'emit');
    component.getFormControl('method').setValue(component.webhookMethods[0]);
    component.getFormControl('url').setValue('https://example.com');
    component.addHeader();
    component.headerControls[0].get('name')?.setValue('x-token');
    component.headerControls[0].get('value')?.setValue('token-value');
    component.getFormControl('payload').setValue('payload');
    component.getFormControl('proxy').setValue('https://proxy.com');

    // when
    component.onWebhookFormChange();

    // then
    expect(spy).toHaveBeenCalledTimes(1);
    expect(spy).toHaveBeenCalledWith({
      header: [{ name: 'x-token', value: 'token-value' }],
      method: 'GET',
      payload: 'payload',
      proxy: 'https://proxy.com',
      sendFinished: true,
      type: '',
      url: 'https://example.com',
    });
  });

  it('sendFinished should be enabled for triggered events and true by default', () => {
    // given
    const checkbox = fixture.nativeElement.querySelector('[uitestid=edit-webhook-field-sendFinished] input');
    component.eventType = 'triggered';
    fixture.detectChanges();

    // when

    // then
    expect(checkbox.disabled).toEqual(false);
    expect(component.getFormControl('sendFinished').value).toEqual('true');
  });

  it('sendFinished should be disabled for started events and null', () => {
    // given
    const checkbox = fixture.nativeElement.querySelector('[uitestid=edit-webhook-field-sendFinished] input');
    component.eventType = 'started';
    fixture.detectChanges();

    // when

    // then
    expect(checkbox.disabled).toEqual(true);
    expect(component.getFormControl('sendFinished').value).toEqual(null);
  });

  it('sendFinished should be disabled for finished events and null', () => {
    // given
    const checkbox = fixture.nativeElement.querySelector('[uitestid=edit-webhook-field-sendFinished] input');
    component.eventType = 'finished';
    fixture.detectChanges();

    // when

    // then
    expect(checkbox.disabled).toEqual(true);
    expect(component.getFormControl('sendFinished').value).toEqual(null);
  });

  it('sendFinished should be set to true', () => {
    // given
    component.eventType = 'triggered';
    component.webhook = {
      header: [{ name: 'x-token', value: 'token-value' }],
      method: 'GET',
      payload: 'payload',
      proxy: 'https://proxy.com',
      sendFinished: true,
      filter: {
        projects: null,
        services: null,
        stages: null,
      },
      type: '',
      url: 'https://example.com',
    };

    // when

    // then
    expect(component.getFormControl('sendFinished').value).toEqual('true');
  });

  it('sendFinished should be set to false', () => {
    // given
    component.eventType = 'triggered';
    component.webhook = {
      header: [{ name: 'x-token', value: 'token-value' }],
      method: 'GET',
      payload: 'payload',
      proxy: 'https://proxy.com',
      sendFinished: false,
      filter: {
        projects: null,
        services: null,
        stages: null,
      },
      type: '',
      url: 'https://example.com',
    };

    // when

    // then
    expect(component.getFormControl('sendFinished').value).toEqual('false');
  });

  function getAddHeaderButton(): HTMLElement {
    return fixture.nativeElement.querySelector('[uitestid="ktb-webhook-settings-add-header-button"]');
  }
});
