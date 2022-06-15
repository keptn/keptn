import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbDeploymentStageTimelineComponent } from './ktb-deployment-stage-timeline.component';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { FlexLayoutModule } from '@angular/flex-layout';
import { KtbStageBadgeModule } from '../ktb-stage-badge/ktb-stage-badge.module';

@NgModule({
  declarations: [KtbDeploymentStageTimelineComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    FlexLayoutModule,
    KtbStageBadgeModule,
  ],
  exports: [KtbDeploymentStageTimelineComponent],
})
export class KtbDeploymentStageTimelineModule {}
