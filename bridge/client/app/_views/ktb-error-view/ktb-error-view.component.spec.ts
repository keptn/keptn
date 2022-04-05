import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbErrorViewComponent } from './ktb-error-view.component';
import { AppModule } from '../../app.module';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { BehaviorSubject } from 'rxjs';

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
      imports: [AppModule],
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

  it('should set error to INTERNAL if status is not provided', () => {
    fixture.detectChanges();
    expect(component.error).toBe(500);
  });

  it('should set error to INTERNAL if status is 500', () => {
    setStatus('500');
    fixture.detectChanges();
    expect(component.error).toBe(500);
  });

  it('should set error to INTERNAL if status is not a number', () => {
    setStatus('abc');
    fixture.detectChanges();
    expect(component.error).toBe(500);
  });

  it('should set error to INSUFFICIENT_PERMISSION if status is 403', () => {
    setStatus('403');
    fixture.detectChanges();
    expect(component.error).toBe(403);
  });

  function setStatus(status?: string): void {
    queryParamMapSubject.next(convertToParamMap(status !== undefined ? { status } : {}));
  }
});
