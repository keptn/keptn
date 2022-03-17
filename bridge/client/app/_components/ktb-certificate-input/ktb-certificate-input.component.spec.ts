import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbCertificateInputComponent } from './ktb-certificate-input.component';

describe('KtbCertificateInputComponent', () => {
  let component: KtbCertificateInputComponent;
  let fixture: ComponentFixture<KtbCertificateInputComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbCertificateInputComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbCertificateInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
