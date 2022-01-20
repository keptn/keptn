import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbProjectCreateMessageComponent } from './ktb-project-create-message.component';
import { AppModule } from '../../../app.module';

describe('KtbProjectCreateMessageComponent', () => {
  let component: KtbProjectCreateMessageComponent;
  let fixture: ComponentFixture<KtbProjectCreateMessageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
      declarations: [],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbProjectCreateMessageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
