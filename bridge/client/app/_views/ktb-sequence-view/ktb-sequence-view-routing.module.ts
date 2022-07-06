import { RouterModule, Routes } from '@angular/router';
import { KtbSequenceViewComponent } from './ktb-sequence-view.component';
import { NgModule } from '@angular/core';

const routes: Routes = [
  {
    path: '',
    component: KtbSequenceViewComponent,
  },
  { path: ':shkeptncontext', component: KtbSequenceViewComponent },
  { path: ':shkeptncontext/event/:eventId', component: KtbSequenceViewComponent },
  { path: ':shkeptncontext/stage/:stage', component: KtbSequenceViewComponent },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class KtbSequenceViewRoutingModule {}
