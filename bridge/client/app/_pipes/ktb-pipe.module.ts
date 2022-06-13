import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TruncateNumberPipe } from './truncate-number';
import { SanitizeHtmlPipe } from './sanitize-html.pipe';
import { ToType } from './to-type';

@NgModule({
  declarations: [TruncateNumberPipe, SanitizeHtmlPipe, ToType],
  imports: [CommonModule],
  exports: [TruncateNumberPipe, SanitizeHtmlPipe, ToType],
})
export class KtbPipeModule {}
