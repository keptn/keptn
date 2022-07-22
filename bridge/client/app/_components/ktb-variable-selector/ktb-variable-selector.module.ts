import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbTreeListSelectModule } from '../ktb-tree-list-select/ktb-tree-list-select.module';
import { KtbVariableSelectorComponent } from './ktb-variable-selector.component';

@NgModule({
  declarations: [KtbVariableSelectorComponent],
  imports: [
    CommonModule,
    HttpClientModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    KtbTreeListSelectModule,
  ],
  exports: [KtbVariableSelectorComponent],
})
export class KtbVariableSelectorModule {}
