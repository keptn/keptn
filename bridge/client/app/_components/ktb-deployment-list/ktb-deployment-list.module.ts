import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbDeploymentListComponent } from './ktb-deployment-list.component';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';

@NgModule({
  declarations: [KtbDeploymentListComponent],
  imports: [
    CommonModule,
    DtTableModule,
    DtTagModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    KtbPipeModule,
  ],
  exports: [KtbDeploymentListComponent],
})
export class KtbDeploymentListModule {}
