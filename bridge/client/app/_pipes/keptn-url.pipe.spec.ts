import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ApiService } from '../_services/api.service';
import { ApiServiceMock } from '../_services/api.service.mock';
import { DataService } from '../_services/data.service';
import { KeptnUrlPipe } from './keptn-url.pipe';
import { of } from 'rxjs';
import { KeptnInfo } from '../_models/keptn-info';

describe('KeptnUrlPipe', () => {
  let dataService: DataService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [HttpClientTestingModule],
      providers: [{ provide: ApiService, useClass: ApiServiceMock }],
    }).compileComponents();

    dataService = TestBed.inject(DataService);
    // eslint-disable-next-line @typescript-eslint/dot-notation
    KeptnUrlPipe['_version'] = undefined;
  });

  it('should show the correct URL', () => {
    jest
      .spyOn(dataService, 'keptnInfo', 'get')
      .mockReturnValue(of({ bridgeInfo: { bridgeVersion: '0.15.5' } } as KeptnInfo));
    const component = creatComponent();

    expect(component.transform('/operate/')).toBe('https://keptn.sh/docs/0.15.x/operate/');
  });

  it('should fallback to install if version is invalid', () => {
    jest
      .spyOn(dataService, 'keptnInfo', 'get')
      .mockReturnValue(of({ bridgeInfo: { bridgeVersion: 'develop' } } as KeptnInfo));
    const component = creatComponent();

    expect(component.transform('/operate/')).toBe('https://keptn.sh/docs/install/');
  });

  it('should fallback to install if version is undefined', () => {
    jest
      .spyOn(dataService, 'keptnInfo', 'get')
      .mockReturnValue(of({ bridgeInfo: { bridgeVersion: undefined } } as KeptnInfo));
    const component = creatComponent();

    expect(component.transform('/operate/')).toBe('https://keptn.sh/docs/install/');
  });

  function creatComponent(): KeptnUrlPipe {
    return new KeptnUrlPipe(dataService);
  }
});
