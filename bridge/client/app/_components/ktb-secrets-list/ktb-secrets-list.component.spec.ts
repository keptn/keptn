import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSecretsListComponent } from './ktb-secrets-list.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbSecretsListModule } from './ktb-secrets-list.module';
import { RouterTestingModule } from '@angular/router/testing';

describe('KtbUniformViewComponent', () => {
  let component: KtbSecretsListComponent;
  let fixture: ComponentFixture<KtbSecretsListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbSecretsListModule, HttpClientTestingModule, RouterTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSecretsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
