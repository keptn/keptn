import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { KtbDashboardLegacyViewComponent } from './ktb-dashboard-legacy-view.component';

const routes: Routes = [
  {
    path: '',
    component: KtbDashboardLegacyViewComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class KtbDashboardLegacyViewRoutingModule {}
