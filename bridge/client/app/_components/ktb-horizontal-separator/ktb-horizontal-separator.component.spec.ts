import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbHorizontalSeparatorComponent } from './ktb-horizontal-separator.component';

describe('KtbHorizontalSeparatorComponent', () => {
  let component: KtbHorizontalSeparatorComponent;
  let fixture: ComponentFixture<KtbHorizontalSeparatorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbHorizontalSeparatorComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbHorizontalSeparatorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
