import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbProjectTileComponent } from './ktb-project-tile.component';
import { KtbDashboardViewModule } from '../ktb-dashboard-view.module';

describe('KtbProjectTileComponent', () => {
  let component: KtbProjectTileComponent;
  let fixture: ComponentFixture<KtbProjectTileComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbDashboardViewModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbProjectTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
