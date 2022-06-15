import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbDragAndDropDirective } from './ktb-drag-and-drop.directive';

@NgModule({
  declarations: [KtbDragAndDropDirective],
  imports: [CommonModule],
  exports: [KtbDragAndDropDirective],
})
export class KtbDragAndDropModule {}
