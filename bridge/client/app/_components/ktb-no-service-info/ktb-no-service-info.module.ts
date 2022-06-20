import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbNoServiceInfoComponent } from './ktb-no-service-info.component';
import { FlexLayoutModule } from '@angular/flex-layout';
import { RouterModule } from '@angular/router';
import { DtIconModule } from '@dynatrace/barista-components/icon';

@NgModule({
  declarations: [KtbNoServiceInfoComponent],
  imports: [
    CommonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    FlexLayoutModule,
    RouterModule,
  ],
  exports: [KtbNoServiceInfoComponent],
})
export class KtbNoServiceInfoModule {}
