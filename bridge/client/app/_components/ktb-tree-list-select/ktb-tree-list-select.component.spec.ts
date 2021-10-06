import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbTreeListSelectComponent, KtbTreeListSelectDirective } from './ktb-tree-list-select.component';
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

describe('KtbTreeListSelectComponent', () => {
  let directive: KtbTreeListSelectDirective;
  let component: KtbTreeListSelectComponent;
  let fixture: ComponentFixture<KtbTreeListSelectComponent>;
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
      .overrideModule(BrowserDynamicTestingModule, { set: { entryComponents: [KtbTreeListSelectComponent] } })
      .compileComponents();

    fixture = TestBed.createComponent(KtbTreeListSelectComponent);
    component = fixture.componentInstance;
    directive = new KtbTreeListSelectDirective(
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

  it('should create the directive', () => {
    expect(directive).toBeTruthy();
  });

  it('should emit the selected secret', () => {
    // given, when
    const spy = jest.spyOn(component.selected, 'emit');
    component.handleClick(row);

    // then
    expect(spy).toHaveBeenCalledWith(testSecretPath);
  });

  it('should emit the secret from the directive when a secret is selected', () => {
    // given, when
    const spy = jest.spyOn(directive.selected, 'emit');
    directive.show();

    // when
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    directive.contentRef?.instance.handleClick(row);

    // then
    expect(spy).toHaveBeenCalledWith(testSecretPath);
  });

  it('should close the dialog when a secret is selected', () => {
    // given
    const spy = jest.spyOn(directive, 'close');
    directive.show();

    // when
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    directive.contentRef?.instance.handleClick(row);

    // then
    expect(spy).toHaveBeenCalled();
  });

  it('should pass data from directive to component', () => {
    // given
    directive.data = [{ name: 'SecretA', keys: [{ name: 'key1' }], path: testSecretPath }];

    // when
    directive.show();
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    const comp = directive.contentRef?.instance;

    // then
    expect(comp?.dataSource.data).toEqual([{ name: 'SecretA', keys: [{ name: 'key1' }], path: testSecretPath }]);
  });

  it('should close the dialog when the components emits a close event', () => {
    // given
    const spy = jest.spyOn(directive, 'close');
    directive.show();
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore // Ignore private property
    const comp = directive.contentRef?.instance;

    // when
    comp?.closeDialog.emit();

    // then
    expect(spy).toHaveBeenCalled();
  });
});
