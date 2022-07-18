import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { KtbCommonUseCasesViewComponent } from './ktb-common-use-cases-view.component';

const routes: Routes = [
  {
    path: '',
    component: KtbCommonUseCasesViewComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class KtbCommonUseCasesViewRoutingModule {}
