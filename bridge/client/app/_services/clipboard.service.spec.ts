import { TestBed } from '@angular/core/testing';
import { DtToast, DtToastModule } from '@dynatrace/barista-components/toast';
import { ClipboardService } from './clipboard.service';
import { AppModule } from '../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('ClipboardService', () => {
  let service: ClipboardService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
        DtToastModule,
      ],
      providers: [DtToast],
    });
    service = TestBed.inject(ClipboardService);
  });

  it('should return beautified JSON strings', () => {
    expect(service.stringify(1)).toEqual('1');
    expect(service.stringify('foobar')).toEqual('foobar');
    expect(service.stringify({foo: 'bar'})).toEqual('{\n  "foo": "bar"\n}');
    expect(service.stringify([{foo: 'bar'}])).toEqual('[\n  {\n    "foo": "bar"\n  }\n]');
  });
});
