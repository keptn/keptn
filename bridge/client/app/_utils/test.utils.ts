export class TestUtils {
  public static createNewDropEventWithFiles(files: File[]): DragEvent {
    const dataTransfer: DataTransfer = new MockDataTransfer(files);
    const event = new DragEvent('drop');
    Object.defineProperty(event.constructor.prototype, 'dataTransfer', {
      value: dataTransfer
    });
    return event;
  }
}

function MockDataTransfer(files) {
  this.dropEffect = 'all';
  this.effectAllowed = 'all';
  this.items = [];
  this.types = ['Files'];
  this.getData = () => {
    return files;
  };
  this.files = [... files];
}
