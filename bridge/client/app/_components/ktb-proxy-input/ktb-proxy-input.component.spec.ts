import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbProxyInputComponent } from './ktb-proxy-input.component';
import { AppModule } from '../../app.module';

describe('KtbProxyInputComponent', () => {
  let component: KtbProxyInputComponent;
  let fixture: ComponentFixture<KtbProxyInputComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbProxyInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should correctly set input data and emit data', () => {
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');
    const triggerValidationHostSpy = jest.spyOn(component.hostControl, 'markAsDirty');
    const triggerValidationPortSpy = jest.spyOn(component.portControl, 'markAsDirty');
    component.proxy = {
      gitProxyUrl: 'http://0.0.0.0:5000',
      gitProxyInsecure: true,
      gitProxyScheme: 'https',
      gitProxyPassword: 'myPassword',
      gitProxyUser: 'myUser',
    };

    expect(component.hostControl.value).toBe('http://0.0.0.0');
    expect(component.portControl.value).toBe('5000');
    expect(component.isInsecureControl.value).toBe(true);
    expect(component.schemeControl.value).toBe('https');
    expect(component.passwordControl.value).toBe('myPassword');
    expect(component.userControl.value).toBe('myUser');
    expect(triggerValidationHostSpy).toHaveBeenCalled();
    expect(triggerValidationPortSpy).toHaveBeenCalled();
    expect(emitSpy).toHaveBeenCalledWith({
      gitProxyUrl: 'http://0.0.0.0:5000',
      gitProxyInsecure: true,
      gitProxyScheme: 'https',
      gitProxyPassword: 'myPassword',
      gitProxyUser: 'myUser',
    });
  });

  it('should correctly set input data and not emit data', () => {
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');
    component.proxy = {
      gitProxyUrl: '',
      gitProxyInsecure: true,
      gitProxyScheme: 'https',
    };

    expect(component.hostControl.value).toBe('');
    expect(component.portControl.value).toBe('');
    expect(component.isInsecureControl.value).toBe(true);
    expect(component.schemeControl.value).toBe('https');
    expect(component.passwordControl.value).toBe('');
    expect(component.userControl.value).toBe('');
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should correctly set input data and emit data without password and user', () => {
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');
    component.proxy = {
      gitProxyUrl: 'http://0.0.0.0:5000',
      gitProxyInsecure: true,
      gitProxyScheme: 'https',
    };

    expect(component.hostControl.value).toBe('http://0.0.0.0');
    expect(component.portControl.value).toBe('5000');
    expect(component.isInsecureControl.value).toBe(true);
    expect(component.schemeControl.value).toBe('https');
    expect(emitSpy).toHaveBeenCalledWith({
      gitProxyUrl: 'http://0.0.0.0:5000',
      gitProxyInsecure: true,
      gitProxyScheme: 'https',
      gitProxyPassword: '',
      gitProxyUser: '',
    });
  });

  it('should emit data if data is set', () => {
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');
    component.hostControl.setValue('https://myHost.com');
    component.portControl.setValue('5000');
    component.passwordControl.setValue('myPassword');
    component.userControl.setValue('myUser');

    component.proxyChanged();

    expect(emitSpy).toHaveBeenCalledWith({
      gitProxyUrl: 'https://myHost.com:5000',
      gitProxyInsecure: false,
      gitProxyScheme: 'https',
      gitProxyPassword: 'myPassword',
      gitProxyUser: 'myUser',
    });
  });

  it('should emit data if data is set without user or password', () => {
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');
    component.hostControl.setValue('https://myHost.com');
    component.portControl.setValue('5000');

    component.proxyChanged();

    expect(emitSpy).toHaveBeenCalledWith({
      gitProxyUrl: 'https://myHost.com:5000',
      gitProxyInsecure: false,
      gitProxyScheme: 'https',
      gitProxyPassword: '',
      gitProxyUser: '',
    });
  });

  it('should emit data if data is set without password', () => {
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');
    component.hostControl.setValue('https://myHost.com');
    component.portControl.setValue('5000');
    component.userControl.setValue('myUser');

    component.proxyChanged();

    expect(emitSpy).toHaveBeenCalledWith({
      gitProxyUrl: 'https://myHost.com:5000',
      gitProxyInsecure: false,
      gitProxyScheme: 'https',
      gitProxyPassword: '',
      gitProxyUser: 'myUser',
    });
  });

  it('should emit data if data is set without user', () => {
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');
    component.hostControl.setValue('https://myHost.com');
    component.portControl.setValue('5000');
    component.passwordControl.setValue('myPassword');

    component.proxyChanged();

    expect(emitSpy).toHaveBeenCalledWith({
      gitProxyUrl: 'https://myHost.com:5000',
      gitProxyInsecure: false,
      gitProxyScheme: 'https',
      gitProxyPassword: 'myPassword',
      gitProxyUser: '',
    });
  });

  it('should emit undefined if port is not set', () => {
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');
    component.hostControl.setValue('https://myHost.com');

    component.proxyChanged();

    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should emit undefined if host is not set', () => {
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');
    component.portControl.setValue('5000');

    component.proxyChanged();

    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should correctly set host and not set port if input data does not contain port', () => {
    for (const host of ['0.0.0.0', 'http://0.0.0.0']) {
      component.proxy = {
        gitProxyUrl: host,
        gitProxyInsecure: true,
        gitProxyScheme: 'https',
        gitProxyPassword: 'myPassword',
        gitProxyUser: 'myUser',
      };

      expect(component.hostControl.value).toBe(host);
      expect(component.portControl.value).toBe('');
    }
  });
});
