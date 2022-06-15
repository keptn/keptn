import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbApprovalItemComponent } from './ktb-approval-item.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbApprovalItemModule } from './ktb-approval-item.module';

describe('KtbEventItemComponent', () => {
  let component: KtbApprovalItemComponent;
  let fixture: ComponentFixture<KtbApprovalItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbApprovalItemModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbApprovalItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
