import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbTreeListSelectComponent, KtbTreeListSelectDirective } from './ktb-tree-list-select.component';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtTreeTableModule } from '@dynatrace/barista-components/tree-table';
import { FlexLayoutModule } from '@angular/flex-layout';

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
