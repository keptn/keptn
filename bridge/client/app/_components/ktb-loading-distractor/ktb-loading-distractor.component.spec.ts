import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AppModule } from '../../app.module';

import { KtbLoadingDistractorComponent } from './ktb-loading-distractor.component';

describe('KtbLoadingDistractorComponent', () => {
  let component: KtbLoadingDistractorComponent;
  let fixture: ComponentFixture<KtbLoadingDistractorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule],
    })
      .compileComponents();

    fixture = TestBed.createComponent(KtbLoadingDistractorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should contain loading spinner with specific class', () => {
    // given
    const element: HTMLElement = fixture.nativeElement;

    // when
    var actual = element.querySelector('ktb-loading-spinner.ktb-loading-spinner');

    // then
    expect(actual).toBeTruthy();
  });

  it('should contain label element with specific class', () => {
    // given
    const element: HTMLElement = fixture.nativeElement;

    // when
    const actual = element.querySelector('.ktb-loading-distractor-label');

    // then
    expect(actual).toBeTruthy();
  });
});
