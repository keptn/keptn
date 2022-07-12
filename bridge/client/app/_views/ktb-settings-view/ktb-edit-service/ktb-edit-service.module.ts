import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtTreeTableModule } from '@dynatrace/barista-components/tree-table';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { KtbLoadingModule } from '../../../_components/ktb-loading/ktb-loading.module';
import { KtbDangerZoneModule } from '../../../_components/ktb-danger-zone/ktb-danger-zone.module';
import { KtbEditServiceComponent } from './ktb-edit-service.component';
import { KtbEditServiceFileListComponent } from './ktb-edit-service-file-list.component';

@NgModule({
  declarations: [KtbEditServiceComponent, KtbEditServiceFileListComponent],
  imports: [
    CommonModule,
    FlexLayoutModule,
    RouterModule,
    DtTreeTableModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    KtbDangerZoneModule,
    KtbLoadingModule,
  ],
  exports: [KtbEditServiceComponent],
})
export class KtbEditServiceModule {}
