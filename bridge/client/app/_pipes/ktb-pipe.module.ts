import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TruncateNumberPipe } from './truncate-number';

@NgModule({
  declarations: [TruncateNumberPipe],
  imports: [CommonModule],
  exports: [TruncateNumberPipe],
})
export class KtbPipeModule {}
