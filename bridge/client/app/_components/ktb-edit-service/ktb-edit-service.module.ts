import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbEditServiceComponent } from './ktb-edit-service.component';
import { KtbEditServiceFileListComponent } from './ktb-edit-service-file-list.component';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { RouterModule } from '@angular/router';
import { FlexLayoutModule } from '@angular/flex-layout';

@NgModule({
  declarations: [KtbEditServiceComponent, KtbEditServiceFileListComponent],
  imports: [CommonModule, KtbLoadingModule, RouterModule, FlexLayoutModule],
  exports: [KtbEditServiceComponent],
})
export class KtbEditServiceModule {}
