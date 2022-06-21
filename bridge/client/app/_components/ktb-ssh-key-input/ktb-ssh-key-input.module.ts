import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtFormFieldModule } from '@dynatrace/barista-components/form-field';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { KtbDragAndDropModule } from '../../_directives/ktb-drag-and-drop/ktb-drag-and-drop.module';
import { KtbSshKeyInputComponent } from './ktb-ssh-key-input.component';

@NgModule({
  declarations: [KtbSshKeyInputComponent],
  imports: [CommonModule, DtFormFieldModule, KtbDragAndDropModule, ReactiveFormsModule, DtInputModule, DtButtonModule],
  exports: [KtbSshKeyInputComponent],
})
export class KtbSshKeyInputModule {}
