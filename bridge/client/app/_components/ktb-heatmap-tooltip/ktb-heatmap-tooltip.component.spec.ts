import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbHeatmapTooltipComponent } from './ktb-heatmap-tooltip.component';

describe('KtbHeatmapTooltipComponent', () => {
  let component: KtbHeatmapTooltipComponent;
  let fixture: ComponentFixture<KtbHeatmapTooltipComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KtbHeatmapTooltipComponent],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbHeatmapTooltipComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
