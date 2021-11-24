import { ComponentFixture, tick } from '@angular/core/testing';
import { Trace } from '../_models/trace';

export class TestUtils {
  public static createNewDropEventWithFiles(files: File[]): DragEvent {
    const dataTransfer: DataTransfer = MockDataTransfer(files);
    const event = new DragEvent('drop');
    Object.defineProperty(event.constructor.prototype, 'dataTransfer', {
      value: dataTransfer,
    });
    Object.defineProperty(event.constructor.prototype, 'preventDefault', {
      value: () => {
        return;
      },
    });
    Object.defineProperty(event.constructor.prototype, 'stopPropagation', {
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

  public static mapTraces(input: unknown): Trace[] {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    return input.map((data: unknown): Trace => {
      tracesMapper(data);
      return Trace.fromJSON(data);

      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      function tracesMapper(trace: any): void {
        trace.traces?.forEach((t: Trace) => {
          tracesMapper(t);
        });
        if (trace.traces) {
          trace.traces = Trace.traceMapper(trace.traces);
        }
      }
    });
  }
}

function MockDataTransfer(files: File[]): DataTransfer {
  /* eslint-disable @typescript-eslint/ban-ts-comment */
  return {
    dropEffect: 'none',
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
  /* eslint-enable @typescript-eslint/ban-ts-comment */
}
