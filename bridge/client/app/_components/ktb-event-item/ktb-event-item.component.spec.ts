import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbEventItemComponent } from './ktb-event-item.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbEventItemModule } from './ktb-event-item.module';

describe('KtbEventItemComponent', () => {
  let component: KtbEventItemComponent;
  let fixture: ComponentFixture<KtbEventItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbEventItemModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbEventItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
