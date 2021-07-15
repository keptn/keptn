import {KtbDragAndDropDirective} from './ktb-drag-and-drop.directive';
import {TestUtils} from '../_utils/test.utils';

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
    const event = TestUtils.createNewDropEventWithFiles([new File(['test'], 'test1.yaml'), new File(['test'], 'test2.yaml')]);

    // when
    const emitSpy = spyOn(directive.onError, 'emit');
    directive.onDrop(event);

    // then
    expect(emitSpy).toHaveBeenCalled();
    expect(emitSpy.calls.mostRecent().args[0]).toEqual('Please select only one file');
  });

  it('should allow more than one file when multiple is true', () => {
    // given
    directive.multiple = true;
    const event = TestUtils.createNewDropEventWithFiles([new File(['test'], 'test1.yaml'), new File(['test'], 'test2.yaml')]);

    // when
    const emitSpy = spyOn(directive.onDropped, 'emit');
    directive.onDrop(event);

    // then
    expect(emitSpy).toHaveBeenCalled();
    expect(emitSpy.calls.mostRecent().args[0].length).toEqual(2);
    expect(emitSpy.calls.mostRecent().args[0][0].name).toEqual('test1.yaml');
  });

  it('should allow only files and not directories', () => {
    // given
    const file = new File(['test'], 'test-folder', {type: ''});
    Object.defineProperty(file.constructor.prototype, 'size', {
      value: 4096
    });
    const event = TestUtils.createNewDropEventWithFiles([file]);

    // when
    const emitSpy = spyOn(directive.onError, 'emit');
    directive.onDrop(event);

    // then
    expect(emitSpy).toHaveBeenCalled();
    expect(emitSpy.calls.mostRecent().args[0]).toEqual('Please select only files');
  });

  it('should allow only a predefined set of file extensions when given', () => {
    // given
    const allowedExtensions = ['yaml', 'yml'];
    directive.allowedExtensions = allowedExtensions;
    const event = TestUtils.createNewDropEventWithFiles([new File(['test'], 'test.png', {type: 'image/png'})]);

    // when
    const emitSpy = spyOn(directive.onError, 'emit');
    directive.onDrop(event);

    // then
    expect(emitSpy).toHaveBeenCalled();
    expect(emitSpy.calls.mostRecent().args[0]).toEqual(`Only ${allowedExtensions.join(', ')} files allowed`);
  });

  it('should should allow only a predefined set of file extensions also for multiple files', () => {
    // given
    const allowedExtensions = ['yaml', 'yml'];
    directive.allowedExtensions = allowedExtensions;
    directive.multiple = true;
    const event = TestUtils.createNewDropEventWithFiles([new File(['test'], 'test1.pdf', {type: 'document/pdf'}), new File(['test'], 'test2.png', {type: 'image/png'})]);

    // when
    const emitSpy = spyOn(directive.onError, 'emit');
    directive.onDrop(event);

    // then
    expect(emitSpy).toHaveBeenCalled();
    expect(emitSpy.calls.mostRecent().args[0]).toEqual(`Only ${allowedExtensions.join(', ')} files allowed`);
  });

  it('should allow all file extensions when not set', () => {
    // given
    const event = TestUtils.createNewDropEventWithFiles([new File(['test'], 'test.jpeg')]);

    // when
    const emitSpy = spyOn(directive.onDropped, 'emit');
    directive.onDrop(event);

    // then
    expect(emitSpy).toHaveBeenCalled();
    expect(emitSpy.calls.mostRecent().args[0].length).toEqual(1);
    expect(emitSpy.calls.mostRecent().args[0][0].name).toEqual('test.jpeg');
  });
});
