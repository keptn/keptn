import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSecretsListComponent } from './ktb-secrets-list.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbUniformViewComponent', () => {
  let component: KtbSecretsListComponent;
  let fixture: ComponentFixture<KtbSecretsListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSecretsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
