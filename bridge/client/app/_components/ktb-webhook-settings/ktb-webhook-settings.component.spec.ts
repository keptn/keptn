import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbWebhookSettingsComponent } from './ktb-webhook-settings.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AbstractControl } from '@angular/forms';
import { WebhookConfigMock } from '../../_services/_mockData/webhook-config.mock';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { SecretScope } from '../../../../shared/interfaces/secret';
import { Secret } from '../../_models/secret';

describe('KtbWebhookSettingsComponent', () => {
  const secretPath = 'SecretA.key1';
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
    urlControl.setValue('keptn.sh');
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

  it('should insert the processed string as value to the payload form field', () => {
    // given
    component.getFormControl('url').setValue('');

    // when
    component.setSecret(secretPath, 'url', 0);

    // then
    expect(component.getFormControl('url').value).toEqual(`{{.${secretPath}}}`);
  });

  it('should insert the processed string as value to the url form field at the given position', () => {
    // given
    component.getFormControl('url').setValue('https://example.com?somestringtoinsert');

    // when
    component.setSecret(secretPath, 'url', 30);

    // then
    expect(component.getFormControl('url').value).toEqual(`https://example.com?somestring{{.${secretPath}}}toinsert`);
  });

  it('should insert the processed string as value to the url form field', () => {
    // given
    component.getFormControl('payload').setValue('');

    // when
    component.setSecret(secretPath, 'payload', 0);

    // then
    expect(component.getFormControl('payload').value).toEqual(`{{.${secretPath}}}`);
  });

  it('should insert the processed string as value to the payload form field at the given position', () => {
    // given
    component.getFormControl('payload').setValue('{id: , project: sockshop}');

    // when
    component.setSecret(secretPath, 'payload', 5);

    // then
    expect(component.getFormControl('payload').value).toEqual(`{id: {{.${secretPath}}}, project: sockshop}`);
  });

  it('should insert the processed string as value to the given header field in the form array', () => {
    // given
    component.addHeader('header1', 'value1');
    component.addHeader('header2', '');

    // when
    component.setSecret(secretPath, 'header', 0, 1);

    // then
    expect(component.headerControls[1].get('value')?.value).toEqual(`{{.${secretPath}}}`);
  });

  it('should insert the processed string as value to the header form field at the given position', () => {
    // given
    component.addHeader('header1', 'value1');
    component.addHeader('header2', 'value2');
    // when
    component.setSecret(secretPath, 'header', 5, 1);

    // then
    expect(component.headerControls[1].get('value')?.value).toEqual(`value{{.${secretPath}}}2`);
  });

  it('should map secrets to a tree when set', () => {
    // given, when
    const secrets = [new Secret(), new Secret()];
    secrets[0].name = 'SecretA';
    secrets[0].scope = SecretScope.WEBHOOK;
    secrets[0].keys = ['key1', 'key2', 'key3'];
    secrets[1].name = 'SecretB';
    secrets[1].scope = SecretScope.WEBHOOK;
    secrets[1].keys = ['key1', 'key2', 'key3'];

    // when
    component.secrets = secrets;

    // then
    expect(component.secretDataSource).toEqual(secretDataSource);
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
    expect(spy).toHaveBeenCalledWith({header: [{name: 'x-token', value: 'token-value'}], method: 'GET', payload: 'payload', proxy: 'https://proxy.com', type: '', url: 'https://example.com'});
  });

  function getAddHeaderButton(): HTMLElement {
    return fixture.nativeElement.querySelector('[uitestid="ktb-webhook-settings-add-header-button"]');
  }
});


const secretDataSource = [
  {
    name: 'SecretA',
    keys: [
      {
        name: 'key1',
        path: 'SecretA.key1',
      },
      {
        name: 'key2',
        path: 'SecretA.key2',
      },
      {
        name: 'key3',
        path: 'SecretA.key3',
      },
    ],
  },
  {
    name: 'SecretB',
    keys: [
      {
        name: 'key1',
        path: 'SecretB.key1',
      },
      {
        name: 'key2',
        path: 'SecretB.key2',
      },
      {
        name: 'key3',
        path: 'SecretB.key3',
      },
    ],
  },
];
