import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { KtbLoadingDistractorComponent } from './ktb-loading-distractor.component';
import { KtbLoadingSpinnerComponent } from './ktb-loading-spinner.component';

@NgModule({
  declarations: [KtbLoadingDistractorComponent, KtbLoadingSpinnerComponent],
  imports: [CommonModule, BrowserAnimationsModule],
  exports: [KtbLoadingDistractorComponent, KtbLoadingSpinnerComponent],
})
export class KtbLoadingModule {}
