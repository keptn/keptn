import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { DtExpandablePanelModule } from '@dynatrace/barista-components/expandable-panel';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtShowMoreModule } from '@dynatrace/barista-components/show-more';
import { KtbExpandableTileComponent, KtbExpandableTileHeaderDirective } from './ktb-expandable-tile.component';

@NgModule({
  declarations: [KtbExpandableTileComponent, KtbExpandableTileHeaderDirective],
  imports: [
    CommonModule,
    DtExpandablePanelModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtShowMoreModule,
  ],
  exports: [KtbExpandableTileComponent, KtbExpandableTileHeaderDirective],
})
export class KtbExpandableTileModule {}
