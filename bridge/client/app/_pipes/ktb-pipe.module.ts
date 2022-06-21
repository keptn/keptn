import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { KeptnUrlPipe } from './keptn-url.pipe';
import { SanitizeHtmlPipe } from './sanitize-html.pipe';
import { ToDatePipe } from './to-date.pipe';
import { ToType } from './to-type';
import { TruncateNumberPipe } from './truncate-number';

@NgModule({
  declarations: [TruncateNumberPipe, SanitizeHtmlPipe, ToType, ToDatePipe, KeptnUrlPipe],
  imports: [CommonModule],
  exports: [TruncateNumberPipe, SanitizeHtmlPipe, ToType, ToDatePipe, KeptnUrlPipe],
})
export class KtbPipeModule {}
