import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbStageBadgeComponent } from './ktb-stage-badge.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbStageBadgeModule } from './ktb-stage-badge.module';

describe('KtbStageBadgeComponent', () => {
  let component: KtbStageBadgeComponent;
  let fixture: ComponentFixture<KtbStageBadgeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbStageBadgeModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbStageBadgeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
