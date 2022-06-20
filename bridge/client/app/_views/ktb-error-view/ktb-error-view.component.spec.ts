import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { BehaviorSubject, firstValueFrom } from 'rxjs';
import { KtbErrorViewComponent } from './ktb-error-view.component';
import { KtbErrorViewModule } from './ktb-error-view.module';

describe('KtbErrorViewComponent', () => {
  let component: KtbErrorViewComponent;
  let fixture: ComponentFixture<KtbErrorViewComponent>;
  const queryParamMapSubject = new BehaviorSubject(
    convertToParamMap({
      status: '500',
    })
  );

  beforeEach(async () => {
    setStatus();
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbErrorViewModule, RouterTestingModule],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            queryParamMap: queryParamMapSubject.asObservable(),
          },
        },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbErrorViewComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should set error to INTERNAL if status is not provided', async () => {
    fixture.detectChanges();
    expect(await getError()).toBe(500);
  });

  it('should set error to INTERNAL if status is 500', async () => {
    setStatus('500');
    fixture.detectChanges();
    expect(await getError()).toBe(500);
  });

  it('should set error to INTERNAL if status is not a number', async () => {
    setStatus('abc');
    fixture.detectChanges();
    expect(await getError()).toBe(500);
  });

  it('should set error to INSUFFICIENT_PERMISSION if status is 403', async () => {
    setStatus('403');
    fixture.detectChanges();
    expect(await getError()).toBe(403);
  });

  it('should set error to NOT_ALLOWED if provided by input', async () => {
    setStatus('500');
    component.error = 405;
    fixture.detectChanges();
    expect(await getError()).toBe(405);
  });

  function setStatus(status?: string): void {
    queryParamMapSubject.next(convertToParamMap(status !== undefined ? { status } : {}));
  }

  function getError(): Promise<number> {
    return firstValueFrom(component.error$);
  }
});
