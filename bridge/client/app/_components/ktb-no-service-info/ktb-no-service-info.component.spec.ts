import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbNoServiceInfoComponent } from './ktb-no-service-info.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbNoServiceInfoComponent', () => {
  let component: KtbNoServiceInfoComponent;
  let fixture: ComponentFixture<KtbNoServiceInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbNoServiceInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
