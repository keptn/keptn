import { NgModule } from '@angular/core';
import { KtbProjectViewRoutingCreateModule } from './ktb-project-view-routing-create.module';
import { KtbProjectViewCommonModule } from './ktb-project-view-common.module';

@NgModule({
  imports: [KtbProjectViewCommonModule, KtbProjectViewRoutingCreateModule],
})
export class KtbProjectViewCreateModule {}
