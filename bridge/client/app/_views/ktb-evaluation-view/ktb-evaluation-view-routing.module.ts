import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { KtbEvaluationViewComponent } from './ktb-evaluation-view.component';

const routes: Routes = [
  { path: '', pathMatch: 'full', component: KtbEvaluationViewComponent },
  { path: ':eventselector', component: KtbEvaluationViewComponent },
];

@NgModule({
  declarations: [],
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class KtbEvaluationViewRoutingModule {}
