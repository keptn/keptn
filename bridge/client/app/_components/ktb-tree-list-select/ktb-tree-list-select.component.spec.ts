import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbTreeListSelectComponent } from './ktb-tree-list-select.component';

describe('KtbTreeListSelectComponent', () => {
  let component: KtbTreeListSelectComponent;
  let fixture: ComponentFixture<KtbTreeListSelectComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbTreeListSelectComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbTreeListSelectComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
