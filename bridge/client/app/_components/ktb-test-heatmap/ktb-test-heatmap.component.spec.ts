import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbTestHeatmapComponent } from './ktb-test-heatmap.component';

describe('KtbTestHeatmapComponent', () => {
  let component: KtbTestHeatmapComponent;
  let fixture: ComponentFixture<KtbTestHeatmapComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KtbTestHeatmapComponent],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbTestHeatmapComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
