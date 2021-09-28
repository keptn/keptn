import { TestBed } from '@angular/core/testing';
import { TemplateRef, ViewContainerRef } from '@angular/core';
import { KtbShowHttpLoadingDirective } from './ktb-show-http-loading.directive';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('hideHttpLoadingDirective', () => {
  let directive: KtbShowHttpLoadingDirective;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [KtbShowHttpLoadingDirective, TemplateRef, ViewContainerRef, HttpClientTestingModule],
    });
    directive = TestBed.inject(KtbShowHttpLoadingDirective);
  });

  it('should be an instance', () => {
    expect(directive).toBeTruthy();
  });
});
