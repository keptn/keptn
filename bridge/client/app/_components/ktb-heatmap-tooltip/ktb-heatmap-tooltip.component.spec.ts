import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbHeatmapTooltipComponent } from './ktb-heatmap-tooltip.component';
import { AppModule } from '../../app.module';

describe('KtbHeatmapTooltipComponent', () => {
  let component: KtbHeatmapTooltipComponent;
  let fixture: ComponentFixture<KtbHeatmapTooltipComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbHeatmapTooltipComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
