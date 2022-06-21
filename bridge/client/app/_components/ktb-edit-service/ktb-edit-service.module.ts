import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbEditServiceComponent } from './ktb-edit-service.component';
import { KtbEditServiceFileListComponent } from './ktb-edit-service-file-list.component';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { RouterModule } from '@angular/router';
import { FlexLayoutModule } from '@angular/flex-layout';
import { KtbDangerZoneModule } from '../ktb-danger-zone/ktb-danger-zone.module';
import { DtTreeTableModule } from '@dynatrace/barista-components/tree-table';
import { DtIconModule } from '@dynatrace/barista-components/icon';

@NgModule({
  declarations: [KtbEditServiceComponent, KtbEditServiceFileListComponent],
  imports: [
    CommonModule,
    DtTreeTableModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    FlexLayoutModule,
    KtbDangerZoneModule,
    KtbLoadingModule,
    RouterModule,
  ],
  exports: [KtbEditServiceComponent],
})
export class KtbEditServiceModule {}
