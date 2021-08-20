import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbUniformViewComponent } from './ktb-uniform-view.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbUniformViewComponent', () => {
  let component: KtbUniformViewComponent;
  let fixture: ComponentFixture<KtbUniformViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbUniformViewComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
