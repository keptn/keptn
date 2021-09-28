import { ComponentFixture, tick } from '@angular/core/testing';

export class TestUtils {
  public static createNewDropEventWithFiles(files: File[]): DragEvent {
    const dataTransfer: DataTransfer = MockDataTransfer(files);
    const event = new DragEvent('drop');
    Object.defineProperty(event.constructor.prototype, 'dataTransfer', {
      value: dataTransfer,
    });
    Object.defineProperty(event.constructor.prototype, 'preventDefault', {
      // eslint-disable-next-line arrow-body-style
      value: () => {
        return;
      },
    });
    Object.defineProperty(event.constructor.prototype, 'stopPropagation', {
      // eslint-disable-next-line arrow-body-style
      value: () => {
        return;
      },
    });
    return event;
  }

  public static mockWindowMatchMedia(): void {
    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: jest.fn().mockImplementation((query) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: jest.fn(), // Deprecated
        removeListener: jest.fn(), // Deprecated
        addEventListener: jest.fn(),
        removeEventListener: jest.fn(),
        dispatchEvent: jest.fn(),
      })),
    });
  }

  /**
   * dt-confirmation-dialog has a really strange behavior. This function is used to update the dialog according to its state
   */
  public static updateDialog(fixture: ComponentFixture<unknown>): void {
    fixture.detectChanges();
    tick();
    fixture.detectChanges();
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
    getData(): File[] {
      return files;
    },
    // @ts-ignore
    files: [...files],
  };
}
