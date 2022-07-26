import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { KtbIntegrationViewComponent } from './ktb-integration-view.component';
import { KtbModifyUniformSubscriptionComponent } from './ktb-modify-uniform-subscription/ktb-modify-uniform-subscription.component';
import { PendingChangesGuard } from '../../../_guards/pending-changes.guard';

const routes: Routes = [
  { path: '', component: KtbIntegrationViewComponent },
  { path: ':integrationId', component: KtbIntegrationViewComponent },
  {
    path: ':integrationId/subscriptions/add',
    component: KtbModifyUniformSubscriptionComponent,
    canDeactivate: [PendingChangesGuard],
  },
  {
    path: ':integrationId/subscriptions/:subscriptionId/edit',
    component: KtbModifyUniformSubscriptionComponent,
    canDeactivate: [PendingChangesGuard],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class KtbIntegrationViewRoutingModule {}
