import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbExpandableTileComponent, KtbExpandableTileHeaderDirective } from './ktb-expandable-tile.component';
import { DtExpandablePanelModule } from '@dynatrace/barista-components/expandable-panel';
import { DtShowMoreModule } from '@dynatrace/barista-components/show-more';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

@NgModule({
  declarations: [KtbExpandableTileComponent, KtbExpandableTileHeaderDirective],
  imports: [
    CommonModule,
    BrowserAnimationsModule,
    DtExpandablePanelModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtShowMoreModule,
  ],
  exports: [KtbExpandableTileComponent, KtbExpandableTileHeaderDirective],
})
export class KtbExpandableTileModule {}
