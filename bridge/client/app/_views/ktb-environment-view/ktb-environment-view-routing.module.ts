import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { KtbEnvironmentViewComponent } from './ktb-environment-view.component';

const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    component: KtbEnvironmentViewComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class KtbEnvironmentViewRoutingModule {}
