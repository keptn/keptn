import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbRootEventsListComponent } from './ktb-root-events-list.component';

describe('KtbEventsListComponent', () => {
  let component: KtbRootEventsListComponent;
  let fixture: ComponentFixture<KtbRootEventsListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbRootEventsListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbRootEventsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
