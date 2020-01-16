import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbHttpLoadingSpinnerComponent } from './ktb-http-loading-spinner.component';

describe('HttpLoadingSpinnerComponent', () => {
  let component: KtbHttpLoadingSpinnerComponent;
  let fixture: ComponentFixture<KtbHttpLoadingSpinnerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbHttpLoadingSpinnerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbHttpLoadingSpinnerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
