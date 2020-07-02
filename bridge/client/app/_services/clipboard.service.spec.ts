import { TestBed } from '@angular/core/testing';
import { DtToast, DtToastModule } from '@dynatrace/barista-components/toast';

import { ClipboardService } from './clipboard.service';

describe('ClipboardService', () => {
  let service: ClipboardService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [DtToastModule],
      providers: [DtToast],
    });
    service = TestBed.inject(ClipboardService);
  });

  it('should return beautified JSON strings', () => {
    expect(service.stringify(1))
      .toBe('1');

    expect(service.stringify('foobar'))
      .toBe('foobar');

    expect(service.stringify({ foo: 'bar' }))
      .toBe('{\n  "foo": "bar"\n}');

    expect(service.stringify([{ foo: 'bar' }]))
      .toBe('[\n  {\n    "foo": "bar"\n  }\n]');
  });
});
