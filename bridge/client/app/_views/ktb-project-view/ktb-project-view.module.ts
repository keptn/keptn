import { NgModule } from '@angular/core';
import { KtbProjectViewCommonModule } from './ktb-project-view-common.module';
import { KtbProjectViewRoutingModule } from './ktb-project-view-routing.module';

@NgModule({
  imports: [KtbProjectViewCommonModule, KtbProjectViewRoutingModule],
})
export class KtbProjectViewModule {}
