import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { KtbErrorViewComponent } from './ktb-error-view.component';
import { DtEmptyStateModule } from '@dynatrace/barista-components/empty-state';

@NgModule({
  declarations: [KtbErrorViewComponent],
  imports: [CommonModule, DtEmptyStateModule, RouterModule],
  exports: [KtbErrorViewComponent],
})
export class KtbErrorViewModule {}
