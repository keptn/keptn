import { TestBed } from '@angular/core/testing';
import { DtToast, DtToastModule } from '@dynatrace/barista-components/toast';
import { ClipboardService } from './clipboard.service';
import { AppModule } from '../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { Clipboard } from '@angular/cdk/clipboard';

describe('ClipboardService', () => {
  let service: ClipboardService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule, DtToastModule],
      providers: [DtToast],
    });
    service = TestBed.inject(ClipboardService);
  });

  it('should beautify JSON and copy to clipboard', () => {
    // given
    const copySpy = jest.spyOn(TestBed.inject(Clipboard), 'copy');
    const toastSpy = jest.spyOn(TestBed.inject(DtToast), 'create');

    // when
    service.copy({ foo: 'bar' }, 'myLabel');

    // then
    expect(copySpy).toHaveBeenCalledWith('{\n  "foo": "bar"\n}');
    expect(toastSpy).toHaveBeenCalledWith('Copied myLabel to clipboard');
  });

  it('should return beautified JSON strings', () => {
    expect(service.stringify(1)).toEqual('1');
    expect(service.stringify(undefined)).toEqual('');
    expect(service.stringify(null)).toEqual('');
    expect(service.stringify('foobar')).toEqual('foobar');
    expect(service.stringify({ foo: 'bar' })).toEqual('{\n  "foo": "bar"\n}');
    expect(service.stringify([{ foo: 'bar' }])).toEqual('[\n  {\n    "foo": "bar"\n  }\n]');
  });
});
