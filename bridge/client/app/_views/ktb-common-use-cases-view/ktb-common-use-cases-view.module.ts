import { NgModule } from '@angular/core';
import { KtbCommonUseCasesViewComponent } from './ktb-common-use-cases-view.component';
import { CommonModule } from '@angular/common';
import { KtbCommonUseCasesViewRoutingModule } from './ktb-common-use-cases-view-routing.module';
import { KtbPipeModule } from '../../_pipes/ktb-pipe.module';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DtExpandablePanelModule } from '@dynatrace/barista-components/expandable-panel';
import { DtShowMoreModule } from '@dynatrace/barista-components/show-more';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { KtbMarkdownComponent } from './ktb-markdown/ktb-markdown.component';

@NgModule({
  declarations: [KtbCommonUseCasesViewComponent, KtbMarkdownComponent],
  imports: [
    CommonModule,
    FlexLayoutModule,
    KtbCommonUseCasesViewRoutingModule,
    KtbPipeModule,
    DtButtonModule,
    DtExpandablePanelModule,
    DtShowMoreModule,
  ],
})
export class KtbCommonUseCasesViewModule {}
