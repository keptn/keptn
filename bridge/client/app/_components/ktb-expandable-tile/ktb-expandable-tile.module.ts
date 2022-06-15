import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbExpandableTileComponent, KtbExpandableTileHeaderDirective } from './ktb-expandable-tile.component';
import { DtExpandablePanelModule } from '@dynatrace/barista-components/expandable-panel';
import { DtShowMoreModule } from '@dynatrace/barista-components/show-more';

@NgModule({
  declarations: [KtbExpandableTileComponent, KtbExpandableTileHeaderDirective],
  imports: [CommonModule, DtExpandablePanelModule, DtShowMoreModule],
  exports: [KtbExpandableTileComponent, KtbExpandableTileHeaderDirective],
})
export class KtbExpandableTileModule {}
