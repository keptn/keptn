export class TestUtils {
  public static createNewDropEventWithFiles(files: File[]): DragEvent {
    const dataTransfer: DataTransfer = MockDataTransfer(files);
    const event = new DragEvent('drop');
    Object.defineProperty(event.constructor.prototype, 'dataTransfer', {
      value: dataTransfer
    });
    return event;
  }
}

function MockDataTransfer(files: File[]): DataTransfer {
  return {
    // @ts-ignore
    dropEffect: 'all',
    effectAllowed: 'all',
    // @ts-ignore
    items: [],
    types: ['Files'],
    // @ts-ignore
    getData() {
      return files;
    },
    // @ts-ignore
    files: [...files]
  };
}
