import { KtbHideHttpLoadingDirective } from './ktb-hide-http-loading.directive';
import { TestBed } from '@angular/core/testing';
import { TemplateRef, ViewContainerRef } from '@angular/core';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('hideHttpLoadingDirective', () => {
  let directive: KtbHideHttpLoadingDirective;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        KtbHideHttpLoadingDirective,
        TemplateRef,
        ViewContainerRef,
        HttpClientTestingModule
      ],
    });
    directive = TestBed.inject(KtbHideHttpLoadingDirective);
  });

  it('should be an instance', () => {
    expect(directive).toBeTruthy();
  });
});
