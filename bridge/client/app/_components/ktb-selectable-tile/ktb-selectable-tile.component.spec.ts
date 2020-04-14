import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSelectableTileComponent } from './ktb-selectable-tile.component';

describe('KtbSelectableTileComponent', () => {
  let component: KtbSelectableTileComponent;
  let fixture: ComponentFixture<KtbSelectableTileComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbSelectableTileComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSelectableTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
