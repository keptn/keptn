import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtTreeTableModule } from '@dynatrace/barista-components/tree-table';
import { KtbTreeListSelectComponent, KtbTreeListSelectDirective } from './ktb-tree-list-select.component';

@NgModule({
  declarations: [KtbTreeListSelectDirective, KtbTreeListSelectComponent],
  imports: [
    CommonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtOverlayModule,
    DtTreeTableModule,
    FlexLayoutModule,
  ],
  exports: [KtbTreeListSelectDirective, KtbTreeListSelectComponent],
})
export class KtbTreeListSelectModule {}
