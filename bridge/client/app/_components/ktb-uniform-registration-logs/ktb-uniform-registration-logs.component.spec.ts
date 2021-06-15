import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbUniformRegistrationLogsComponent } from './ktb-uniform-registration-logs.component';

describe('KtbUniformRegistrationLogsComponent', () => {
  let component: KtbUniformRegistrationLogsComponent;
  let fixture: ComponentFixture<KtbUniformRegistrationLogsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbUniformRegistrationLogsComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbUniformRegistrationLogsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
