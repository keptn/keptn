import { FormUtils } from './form.utils';

describe('FormUtils', () => {
  it('should remove every tab, space and newline', () => {
    expect(FormUtils.removeWhitespaces('  a \t\n  b  \n  c \t d  ')).toBe('abcd');
  });
  it('should not remove other characters than tab, space and newline', () => {
    expect(FormUtils.removeWhitespaces('my-r\\epo/my-app_ng:v1.2.0a')).toBe('my-r\\epo/my-app_ng:v1.2.0a');
  });
  it('should not remove anything for a simple input value', () => {
    expect(FormUtils.removeWhitespaces('abcd:1234')).toBe('abcd:1234');
  });
});
