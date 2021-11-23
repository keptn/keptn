import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSecretSelectorComponent } from './ktb-secret-selector.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { Overlay, OverlayPositionBuilder } from '@angular/cdk/overlay';
import { ElementRef } from '@angular/core';
import { Router } from '@angular/router';
import { BrowserDynamicTestingModule } from '@angular/platform-browser-dynamic/testing';

export class MockElementRef extends ElementRef {
  nativeElement = {};

  constructor() {
    super(null);
  }
}

describe('KtbSecretSelectorComponent', () => {
  let component: KtbSecretSelectorComponent;
  let fixture: ComponentFixture<KtbSecretSelectorComponent>;
  const testSecretPath = 'SecretA.key1';
  const row = {
    name: 'key1',
    level: 1,
    path: testSecretPath,
    expandable: false,
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
      providers: [{ provide: ElementRef, useClass: MockElementRef }],
    })
      .overrideModule(BrowserDynamicTestingModule, { set: { entryComponents: [KtbSecretSelectorComponent] } })
      .compileComponents();

    fixture = TestBed.createComponent(KtbSecretSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create the component', () => {
    expect(component).toBeTruthy();
  });
});
