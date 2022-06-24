import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbSelectableTileComponent, KtbSelectableTileHeaderDirective } from './ktb-selectable-tile.component';

@NgModule({
  declarations: [KtbSelectableTileComponent, KtbSelectableTileHeaderDirective],
  imports: [CommonModule],
  exports: [KtbSelectableTileComponent, KtbSelectableTileHeaderDirective],
})
export class KtbSelectableTileModule {}
