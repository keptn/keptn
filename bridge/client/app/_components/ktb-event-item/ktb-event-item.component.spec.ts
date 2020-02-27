import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEventItemComponent } from './ktb-event-item.component';

describe('KtbEventItemComponent', () => {
  let component: KtbEventItemComponent;
  let fixture: ComponentFixture<KtbEventItemComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbEventItemComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbEventItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
