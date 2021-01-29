import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEnvironmentViewComponent } from './ktb-environment-view.component';

describe('KtbEnvironmentComponent', () => {
  let component: KtbEnvironmentViewComponent;
  let fixture: ComponentFixture<KtbEnvironmentViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbEnvironmentViewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbEnvironmentViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
