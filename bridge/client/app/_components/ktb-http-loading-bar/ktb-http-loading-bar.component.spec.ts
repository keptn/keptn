import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbHttpLoadingBarComponent } from './ktb-http-loading-bar.component';

describe('HttpLoadingBarComponent', () => {
  let component: KtbHttpLoadingBarComponent;
  let fixture: ComponentFixture<KtbHttpLoadingBarComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbHttpLoadingBarComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbHttpLoadingBarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
