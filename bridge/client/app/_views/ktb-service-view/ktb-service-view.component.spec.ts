import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbServiceViewComponent } from './ktb-service-view.component';
import { AppModule } from '../../app.module';

describe('KtbEventsListComponent', () => {
  let component: KtbServiceViewComponent;
  let fixture: ComponentFixture<KtbServiceViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbServiceViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
