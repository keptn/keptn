import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbHeatmapComponent } from './ktb-heatmap.component';

describe('KtbHeatmapComponent', () => {
  let component: KtbHeatmapComponent;
  let fixture: ComponentFixture<KtbHeatmapComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KtbHeatmapComponent],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbHeatmapComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
