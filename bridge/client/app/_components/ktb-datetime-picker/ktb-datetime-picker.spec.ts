import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbDatetimePickerComponent, KtbDatetimePickerDirective } from './ktb-datetime-picker.component';
import { AppModule } from '../../app.module';
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

describe('KtbDatetimePickerComponent', () => {
  let directive: KtbDatetimePickerDirective;
  let component: KtbDatetimePickerComponent;
  let fixture: ComponentFixture<KtbDatetimePickerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
      providers: [{ provide: ElementRef, useClass: MockElementRef }],
    })
      .overrideModule(BrowserDynamicTestingModule, { set: { entryComponents: [KtbDatetimePickerComponent] } })
      .compileComponents();

    fixture = TestBed.createComponent(KtbDatetimePickerComponent);
    component = fixture.componentInstance;
    directive = new KtbDatetimePickerDirective(
      TestBed.inject(Overlay),
      TestBed.inject(OverlayPositionBuilder),
      TestBed.inject(ElementRef),
      TestBed.inject(Router)
    );
    directive.ngOnInit();
    fixture.detectChanges();
  });

  it('should create the component', () => {
    expect(component).toBeTruthy();
  });
});
