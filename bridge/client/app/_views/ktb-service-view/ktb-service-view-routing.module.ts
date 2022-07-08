import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { KtbServiceViewComponent } from './ktb-service-view.component';

const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    component: KtbServiceViewComponent,
  },
  { path: ':serviceName', component: KtbServiceViewComponent },
  { path: ':serviceName/context/:shkeptncontext', component: KtbServiceViewComponent },
  { path: ':serviceName/context/:shkeptncontext/stage/:stage', component: KtbServiceViewComponent },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class KtbServiceViewRoutingModule {}
