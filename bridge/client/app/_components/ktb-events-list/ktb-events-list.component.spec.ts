import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEventsListComponent } from './ktb-events-list.component';

describe('KtbEventsListComponent', () => {
  let component: KtbEventsListComponent;
  let fixture: ComponentFixture<KtbEventsListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbEventsListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbEventsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
