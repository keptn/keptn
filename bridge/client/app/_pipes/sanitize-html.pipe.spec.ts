import { TestBed } from '@angular/core/testing';
import { SanitizeHtmlPipe } from './sanitize-html.pipe';
import { DomSanitizer } from '@angular/platform-browser';

describe('SanitizeHtml', () => {
  let sanitizer: DomSanitizer;
  let pipe: SanitizeHtmlPipe;

  beforeEach(async () => {
    sanitizer = TestBed.inject(DomSanitizer);
    pipe = new SanitizeHtmlPipe(sanitizer);
  });

  it('sanitize simple html', () => {
    expect(pipe.transform('<p>Some text</p>')).toBe('<p>Some text</p>');
  });

  it('sanitize html with some script', () => {
    expect(pipe.transform('<p>Some text</p><script>alert("hello world")</script>')).toBe('<p>Some text</p>');
  });

  it('sanitize html with onerror', () => {
    expect(pipe.transform('<p>Some text</p><img src="not-available.png" onerror=alert("hello world") />')).toBe(
      '<p>Some text</p><img src="not-available.png">'
    );
  });
});
