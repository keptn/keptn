import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbEvaluationViewRoutingModule } from './ktb-evaluation-view-routing.module';
import { KtbEvaluationViewComponent } from './ktb-evaluation-view.component';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtInfoGroupModule } from '@dynatrace/barista-components/info-group';
import { FlexModule } from '@angular/flex-layout';
import { DtEmptyStateModule } from '@dynatrace/barista-components/empty-state';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { KtbLoadingModule } from '../../_components/ktb-loading/ktb-loading.module';
import { KtbEventItemModule } from '../../_components/ktb-event-item/ktb-event-item.module';
import { KtbEvaluationDetailsModule } from '../../_components/ktb-evaluation-details/ktb-evaluation-details.module';

@NgModule({
  declarations: [KtbEvaluationViewComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtEmptyStateModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    DtInfoGroupModule,
    FlexModule,
    KtbEvaluationDetailsModule,
    KtbEvaluationViewRoutingModule,
    KtbEventItemModule,
    KtbLoadingModule,
    KtbPipeModule,
  ],
})
export class KtbEvaluationViewModule {}
