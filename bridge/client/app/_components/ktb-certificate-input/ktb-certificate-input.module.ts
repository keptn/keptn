import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbCertificateInputComponent } from './ktb-certificate-input.component';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { ReactiveFormsModule } from '@angular/forms';
import { KtbLoadingModule } from '../ktb-loading/ktb-loading.module';
import { KtbDragAndDropModule } from '../../_directives/ktb-drag-and-drop/ktb-drag-and-drop.module';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtInputModule } from '@dynatrace/barista-components/input';

@NgModule({
  declarations: [KtbCertificateInputComponent],
  imports: [
    CommonModule,
    DtButtonModule,
    DtFormFieldModule,
    DtInputModule,
    KtbDragAndDropModule,
    KtbLoadingModule,
    ReactiveFormsModule,
  ],
  exports: [KtbCertificateInputComponent],
})
export class KtbCertificateInputModule {}
