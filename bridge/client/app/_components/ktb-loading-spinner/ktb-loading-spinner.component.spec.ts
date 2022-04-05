import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbLoadingSpinnerComponent } from './ktb-loading-spinner.component';

describe('KtbLoadingSpinnerComponent', () => {
  let component: KtbLoadingSpinnerComponent;
  let fixture: ComponentFixture<KtbLoadingSpinnerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KtbLoadingSpinnerComponent],
    })
      .compileComponents();

    fixture = TestBed.createComponent(KtbLoadingSpinnerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should contain SVG element with specific class', () => {
    // given
    const element: HTMLElement = fixture.nativeElement;

    // when
    const actual = element.querySelector('svg.ktb-loading-spinner-svg');

    // then
    expect(actual).toBeTruthy();
  });

  it('should contain spinner path inside SVG element', () => {
    // given
    const element: HTMLElement = fixture.nativeElement;

    // when
    const actual = element.querySelector('svg .ktb-loading-spinner-path');

    // then
    expect(actual).toBeTruthy();
  });
});
