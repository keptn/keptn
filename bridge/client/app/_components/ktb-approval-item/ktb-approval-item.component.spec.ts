import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbApprovalItemComponent } from './ktb-approval-item.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbEventItemComponent', () => {
  let component: KtbApprovalItemComponent;
  let fixture: ComponentFixture<KtbApprovalItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbApprovalItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
