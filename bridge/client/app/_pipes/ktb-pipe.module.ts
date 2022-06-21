import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { KeptnUrlPipe } from './keptn-url.pipe';
import { SanitizeHtmlPipe } from './sanitize-html.pipe';
import { ToDatePipe } from './to-date.pipe';
import { ToType } from './to-type';
import { TruncateNumberPipe } from './truncate-number';
import { ArrayToStringPipe } from './array-to-string';

@NgModule({
  declarations: [TruncateNumberPipe, SanitizeHtmlPipe, ToType, ToDatePipe, KeptnUrlPipe, ArrayToStringPipe],
  imports: [CommonModule],
  exports: [TruncateNumberPipe, SanitizeHtmlPipe, ToType, ToDatePipe, KeptnUrlPipe, ArrayToStringPipe],
})
export class KtbPipeModule {}
