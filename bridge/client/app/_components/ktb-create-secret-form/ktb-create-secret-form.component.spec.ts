import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbCreateSecretFormComponent } from './ktb-create-secret-form.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbCreateSecretFormComponent', () => {
  let component: KtbCreateSecretFormComponent;
  let fixture: ComponentFixture<KtbCreateSecretFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbCreateSecretFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
