import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbVariableSelectorComponent } from './ktb-variable-selector.component';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbTreeListSelectModule } from '../ktb-tree-list-select/ktb-tree-list-select.module';

@NgModule({
  declarations: [KtbVariableSelectorComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    KtbTreeListSelectModule,
  ],
  exports: [KtbVariableSelectorComponent],
})
export class KtbVariableSelectorModule {}
