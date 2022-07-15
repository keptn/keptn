import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbChartComponent } from './ktb-chart.component';
import { KtbChartModule } from './ktb-chart.module';

describe('KtbChartComponent', () => {
  let component: KtbChartComponent;
  let fixture: ComponentFixture<KtbChartComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbChartModule, HttpClientTestingModule],
      providers: [],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbChartComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
