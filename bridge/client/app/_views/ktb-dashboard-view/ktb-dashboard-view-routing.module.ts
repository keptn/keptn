import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { KtbDashboardViewComponent } from './ktb-dashboard-view.component';

const routes: Routes = [
  {
    path: '',
    component: KtbDashboardViewComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class KtbDashboardViewRoutingModule {}
