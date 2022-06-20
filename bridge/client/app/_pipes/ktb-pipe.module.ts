import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { SanitizeHtmlPipe } from './sanitize-html.pipe';
import { ToDatePipe } from './to-date.pipe';
import { ToType } from './to-type';
import { TruncateNumberPipe } from './truncate-number';

@NgModule({
  declarations: [TruncateNumberPipe, SanitizeHtmlPipe, ToType, ToDatePipe],
  imports: [CommonModule],
  exports: [TruncateNumberPipe, SanitizeHtmlPipe, ToType, ToDatePipe],
})
export class KtbPipeModule {}
