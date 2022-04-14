import { ComponentFixture, TestBed, tick } from '@angular/core/testing';
import { Trace } from '../_models/trace';
import { ApiService } from '../_services/api.service';
import { BridgeInfoResponseMock } from '../_services/_mockData/api-responses/bridgeInfo-response.mock';
import { of } from 'rxjs';
import { DataService } from '../_services/data.service';

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

  public static createFileList(content: string): FileList {
    const file: File = {
      name: '',
      lastModified: 0,
      type: 'text/plain',
      text(): Promise<string> {
        return Promise.resolve(content);
      },
      size: 0,
      slice(): Blob {
        return new Blob();
      },
      arrayBuffer(): Promise<ArrayBuffer> {
        return Promise.resolve(new ArrayBuffer(0));
      },
      stream(): ReadableStream<unknown> {
        return new ReadableStream<unknown>();
      },
    };

    return {
      0: file,
      length: 1,
      item(): File {
        return file;
      },
    };
  }

  public static enableResourceService(): void {
    const apiService = TestBed.inject(ApiService);
    const dataService = TestBed.inject(DataService);
    const mock = {
      ...BridgeInfoResponseMock,
      featureFlags: {
        RESOURCE_SERVICE_ENABLED: true,
      },
    };
    jest.spyOn(apiService, 'getKeptnInfo').mockReturnValue(of(mock));
    dataService.loadKeptnInfo();
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
