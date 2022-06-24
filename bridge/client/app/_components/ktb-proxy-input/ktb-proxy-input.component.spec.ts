import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbProxyInputComponent } from './ktb-proxy-input.component';
import { KtbProxyInputModule } from './ktb-proxy-input.module';

describe('KtbProxyInputComponent', () => {
  let component: KtbProxyInputComponent;
  let fixture: ComponentFixture<KtbProxyInputComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        KtbProxyInputModule,
        DtIconModule.forRoot({
          svgIconLocation: `assets/icons/{{name}}.svg`,
        }),
        HttpClientTestingModule,
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbProxyInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should correctly set input data and emit data', () => {
    // given
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');

    // when
    component.proxy = {
      gitProxyUrl: 'http://0.0.0.0:5000',
      gitProxyInsecure: true,
      gitProxyScheme: 'https',
      gitProxyPassword: 'myPassword',
      gitProxyUser: 'myUser',
    };

    // then
    expect(component.hostControl.value).toBe('http://0.0.0.0');
    expect(component.portControl.value).toBe('5000');
    expect(component.isInsecureControl.value).toBe(true);
    expect(component.schemeControl.value).toBe('https');
    expect(component.passwordControl.value).toBe('myPassword');
    expect(component.userControl.value).toBe('myUser');
    expect(emitSpy).not.toHaveBeenCalled();

    // when
    component.proxyChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      gitProxyUrl: 'http://0.0.0.0:5000',
      gitProxyInsecure: true,
      gitProxyScheme: 'https',
      gitProxyPassword: 'myPassword',
      gitProxyUser: 'myUser',
    });
  });

  it('should correctly set input data and not emit data', () => {
    // given
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');

    // when
    component.proxy = {
      gitProxyUrl: '',
      gitProxyInsecure: true,
      gitProxyScheme: 'https',
    };

    // then
    expect(component.hostControl.value).toBe('');
    expect(component.portControl.value).toBe('');
    expect(component.isInsecureControl.value).toBe(true);
    expect(component.schemeControl.value).toBe('https');
    expect(component.passwordControl.value).toBe('');
    expect(component.userControl.value).toBe('');
    expect(emitSpy).not.toHaveBeenCalled();
  });

  it('should correctly set input data and emit data without password and user', () => {
    // given
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');

    // when
    component.proxy = {
      gitProxyUrl: 'http://0.0.0.0:5000',
      gitProxyInsecure: true,
      gitProxyScheme: 'https',
    };

    // then
    expect(component.hostControl.value).toBe('http://0.0.0.0');
    expect(component.portControl.value).toBe('5000');
    expect(component.isInsecureControl.value).toBe(true);
    expect(component.schemeControl.value).toBe('https');
    expect(emitSpy).not.toHaveBeenCalled();

    // when
    component.proxyChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      gitProxyUrl: 'http://0.0.0.0:5000',
      gitProxyInsecure: true,
      gitProxyScheme: 'https',
      gitProxyPassword: '',
      gitProxyUser: '',
    });
  });

  it('should emit data if data is set', () => {
    // given
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');

    // when
    component.hostControl.setValue('https://myHost.com');
    component.portControl.setValue('5000');
    component.passwordControl.setValue('myPassword');
    component.userControl.setValue('myUser');

    component.proxyChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      gitProxyUrl: 'https://myHost.com:5000',
      gitProxyInsecure: false,
      gitProxyScheme: 'https',
      gitProxyPassword: 'myPassword',
      gitProxyUser: 'myUser',
    });
  });

  it('should emit data if data is set without user or password', () => {
    // given
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');

    // when
    component.hostControl.setValue('https://myHost.com');
    component.portControl.setValue('5000');

    component.proxyChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      gitProxyUrl: 'https://myHost.com:5000',
      gitProxyInsecure: false,
      gitProxyScheme: 'https',
      gitProxyPassword: '',
      gitProxyUser: '',
    });
  });

  it('should emit data if data is set without password', () => {
    // given
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');

    // when
    component.hostControl.setValue('https://myHost.com');
    component.portControl.setValue('5000');
    component.userControl.setValue('myUser');

    component.proxyChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      gitProxyUrl: 'https://myHost.com:5000',
      gitProxyInsecure: false,
      gitProxyScheme: 'https',
      gitProxyPassword: '',
      gitProxyUser: 'myUser',
    });
  });

  it('should emit data if data is set without user', () => {
    // given
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');

    // when
    component.hostControl.setValue('https://myHost.com');
    component.portControl.setValue('5000');
    component.passwordControl.setValue('myPassword');

    component.proxyChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      gitProxyUrl: 'https://myHost.com:5000',
      gitProxyInsecure: false,
      gitProxyScheme: 'https',
      gitProxyPassword: 'myPassword',
      gitProxyUser: '',
    });
  });

  it('should emit undefined if port is not set', () => {
    // given
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');

    // when
    component.hostControl.setValue('https://myHost.com');

    component.proxyChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should emit undefined if host is not set', () => {
    // given
    const emitSpy = jest.spyOn(component.proxyChange, 'emit');

    // when
    component.portControl.setValue('5000');

    component.proxyChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should correctly set host and not set port if input data does not contain port', () => {
    for (const host of ['0.0.0.0', 'http://0.0.0.0']) {
      // when
      component.proxy = {
        gitProxyUrl: host,
        gitProxyInsecure: true,
        gitProxyScheme: 'https',
        gitProxyPassword: 'myPassword',
        gitProxyUser: 'myUser',
      };

      // then
      expect(component.hostControl.value).toBe(host);
      expect(component.portControl.value).toBe('');
    }
  });
});
