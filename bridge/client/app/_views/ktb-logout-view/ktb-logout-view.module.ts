import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { KtbLogoutViewComponent } from './ktb-logout-view.component';

const routes: Routes = [
  {
    path: '',
    component: KtbLogoutViewComponent,
  },
];

@NgModule({
  declarations: [KtbLogoutViewComponent],
  imports: [CommonModule, RouterModule.forChild(routes)],
})
export class KtbLogoutViewModule {}
