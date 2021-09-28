import { KtbDragAndDropDirective } from './ktb-drag-and-drop.directive';
import { TestUtils } from '../../_utils/test.utils';

describe('KtbDragAndDropDirective', () => {
  let directive: KtbDragAndDropDirective;

  beforeEach(() => {
    directive = new KtbDragAndDropDirective();
  });

  it('should create an instance', () => {
    expect(directive).toBeTruthy();
  });

  it('should allow only one file when multiple is false', () => {
    // given
    directive.multiple = false;
    const event = TestUtils.createNewDropEventWithFiles([
      new File(['test'], 'test1.yaml', { type: 'some' }),
      new File(['test'], 'test2.yaml', { type: 'some' }),
    ]);

    // when
    const emitSpy = jest.spyOn(directive.dropError, 'emit');
    directive.onDrop(event);

    // then
    expect(emitSpy).toHaveBeenCalled();
    expect(emitSpy).toHaveBeenCalledWith('Please select only one file');
  });

  it('should allow more than one file when multiple is true', () => {
    // given
    directive.multiple = true;
    const event = TestUtils.createNewDropEventWithFiles([
      new File(['test'], 'test1.yaml', { type: 'some' }),
      new File(['test'], 'test2.yaml', { type: 'some' }),
    ]);

    // when
    const emitSpy = jest.spyOn(directive.dropped, 'emit');
    directive.onDrop(event);

    // then
    expect(emitSpy).toHaveBeenCalled();
    expect(emitSpy).toHaveBeenCalledWith(event.dataTransfer?.files);
  });

  it('should allow only files and not directories', () => {
    // given
    const file = new File(['test'], 'test-folder', { type: '' });
    Object.defineProperty(file.constructor.prototype, 'size', {
      value: 4096,
    });
    const event = TestUtils.createNewDropEventWithFiles([file]);

    // when
    const emitSpy = jest.spyOn(directive.dropError, 'emit');
    directive.onDrop(event);

    // then
    expect(emitSpy).toHaveBeenCalled();
    expect(emitSpy).toHaveBeenCalledWith('Please select only files');
  });

  it('should allow only a predefined set of file extensions when given', () => {
    // given
    const allowedExtensions = ['yaml', 'yml'];
    directive.allowedExtensions = allowedExtensions;
    const event = TestUtils.createNewDropEventWithFiles([new File(['test'], 'test.png', { type: 'image/png' })]);

    // when
    const emitSpy = jest.spyOn(directive.dropError, 'emit');
    directive.onDrop(event);

    // then
    expect(emitSpy).toHaveBeenCalled();
    expect(emitSpy).toHaveBeenCalledWith(`Only ${allowedExtensions.join(', ')} files allowed`);
  });

  it('should should allow only a predefined set of file extensions also for multiple files', () => {
    // given
    const allowedExtensions = ['yaml', 'yml'];
    directive.allowedExtensions = allowedExtensions;
    directive.multiple = true;
    const event = TestUtils.createNewDropEventWithFiles([
      new File(['test'], 'test1.pdf', { type: 'document/pdf' }),
      new File(['test'], 'test2.png', { type: 'image/png' }),
    ]);

    // when
    const emitSpy = jest.spyOn(directive.dropError, 'emit');
    directive.onDrop(event);

    // then
    expect(emitSpy).toHaveBeenCalled();
    expect(emitSpy).toHaveBeenCalledWith(`Only ${allowedExtensions.join(', ')} files allowed`);
  });

  it('should allow all file extensions when not set', () => {
    // given
    const event = TestUtils.createNewDropEventWithFiles([new File(['test'], 'test.jpeg', { type: 'image/jpeg' })]);

    // when
    const emitSpy = jest.spyOn(directive.dropped, 'emit');
    directive.onDrop(event);

    // then
    expect(emitSpy).toHaveBeenCalled();
    expect(emitSpy).toHaveBeenCalledWith(event.dataTransfer?.files);
  });
});
