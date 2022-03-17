import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProxyInputComponent } from './ktb-proxy-input.component';

describe('KtbProxyInputComponent', () => {
  let component: KtbProxyInputComponent;
  let fixture: ComponentFixture<KtbProxyInputComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbProxyInputComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbProxyInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
