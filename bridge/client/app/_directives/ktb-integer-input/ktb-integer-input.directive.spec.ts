import { KtbIntegerInputDirective } from './ktb-integer-input.directive';
import { FormControl, NgControl } from '@angular/forms';

describe('KtbIntegerInputDirective', () => {
  let directive: KtbIntegerInputDirective;
  let formControl: FormControl;
  beforeEach(async () => {
    formControl = new FormControl('');
    directive = new KtbIntegerInputDirective({
      control: formControl,
    } as unknown as NgControl);
  });

  it('should create an instance', () => {
    expect(directive).toBeTruthy();
  });

  it('should prevent adding dots on keydown', () => {
    for (const separator of ['.', ',']) {
      // given
      const event = new KeyboardEvent('keydown', {
        key: separator,
      });
      const preventSpy = jest.spyOn(event, 'preventDefault');

      // when
      directive.onKeyDown(event);

      // then
      expect(formControl.value).toBe('');
      expect(preventSpy).toHaveBeenCalled();
    }
  });

  it('should prevent adding minus on keydown', () => {
    // given
    const event = new KeyboardEvent('keydown', {
      key: '-',
    });
    const preventSpy = jest.spyOn(event, 'preventDefault');

    // when
    directive.onKeyDown(event);

    // then
    expect(formControl.value).toBe('');
    expect(preventSpy).toHaveBeenCalled();
  });

  it('should not prevent adding minus on keydown', () => {
    // given
    directive.allowNegative = true;
    const event = new KeyboardEvent('keydown', {
      key: '-',
    });
    const preventSpy = jest.spyOn(event, 'preventDefault');

    // when
    directive.onKeyDown(event);

    // then
    expect(formControl.value).toBe('');
    expect(preventSpy).not.toHaveBeenCalled();
  });

  it('should not prevent entering numbers on keydown', () => {
    for (let i = 0; i < 10; ++i) {
      // given
      const event = new KeyboardEvent('keydown', {
        key: i.toString(),
      });
      const preventSpy = jest.spyOn(event, 'preventDefault');

      // when
      directive.onKeyDown(event);

      // then
      expect(preventSpy).not.toHaveBeenCalled();
    }
  });

  it('should trunc decimals on change/paste', () => {
    // given
    formControl.setValue('123.456');

    // when
    directive.onInput();

    // then
    expect(formControl.value).toBe('123');
  });

  it('should trunc negative on change/paste', () => {
    // given
    formControl.setValue('-123.456');

    // when
    directive.onInput();

    // then
    expect(formControl.value).toBe('123');
  });

  it('should not trunc negative on change/paste', () => {
    // given
    directive.allowNegative = true;
    formControl.setValue('-123.456');

    // when
    directive.onInput();

    // then
    expect(formControl.value).toBe('-123');
  });

  it('should not modify changed/pasted text', () => {
    // given
    formControl.setValue('123');

    // when
    directive.onInput();

    // then
    expect(formControl.value).toBe('123');
  });
});
