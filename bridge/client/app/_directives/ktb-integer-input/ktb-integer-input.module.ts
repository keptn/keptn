import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { KtbIntegerInputDirective } from './ktb-integer-input.directive';

@NgModule({
  declarations: [KtbIntegerInputDirective],
  imports: [CommonModule],
  exports: [KtbIntegerInputDirective],
})
export class KtbIntegerInputModule {}
