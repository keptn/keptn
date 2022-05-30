import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TruncateNumberPipe } from './truncate-number';
import { SanitizeHtmlPipe } from './sanitize-html.pipe';

@NgModule({
  declarations: [TruncateNumberPipe, SanitizeHtmlPipe],
  imports: [CommonModule],
  exports: [TruncateNumberPipe, SanitizeHtmlPipe],
})
export class KtbPipeModule {}
