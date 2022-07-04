import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbLoadingDistractorComponent } from './ktb-loading-distractor.component';
import { KtbLoadingSpinnerComponent } from './ktb-loading-spinner.component';

@NgModule({
  declarations: [KtbLoadingDistractorComponent, KtbLoadingSpinnerComponent],
  imports: [CommonModule],
  exports: [KtbLoadingDistractorComponent, KtbLoadingSpinnerComponent],
})
export class KtbLoadingModule {}
