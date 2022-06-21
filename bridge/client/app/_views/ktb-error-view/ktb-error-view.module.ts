import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { KtbErrorViewComponent } from './ktb-error-view.component';

@NgModule({
  declarations: [KtbErrorViewComponent],
  imports: [CommonModule, RouterModule],
  exports: [KtbErrorViewComponent],
})
export class KtbErrorViewModule {}
