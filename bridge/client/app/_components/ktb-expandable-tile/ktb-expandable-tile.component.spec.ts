import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbExpandableTileComponent } from './ktb-expandable-tile.component';

describe('KtbExpandableTileComponent', () => {
  let component: KtbExpandableTileComponent;
  let fixture: ComponentFixture<KtbExpandableTileComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbExpandableTileComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbExpandableTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
