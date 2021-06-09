import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashboardComponent} from './dashboard/dashboard.component';
import {ProjectBoardComponent} from './project-board/project-board.component';
import {EvaluationBoardComponent} from './evaluation-board/evaluation-board.component';
import {ForwarderGuard} from './_forwarder/forwarder_guard';


const routes: Routes = [
  {path: '', pathMatch: 'full', redirectTo: 'dashboard'},
  {path: 'dashboard', component: DashboardComponent},
  {path: 'project/:projectName', component: ProjectBoardComponent},
  {path: 'project/:projectName/uniform', component: ProjectBoardComponent},
  {path: 'project/:projectName/integration', component: ProjectBoardComponent},
  {path: 'project/:projectName/settings', component: ProjectBoardComponent},
  {path: 'project/:projectName/service', component: ProjectBoardComponent},
  {path: 'project/:projectName/service/:serviceName', component: ProjectBoardComponent},
  {path: 'project/:projectName/service/:serviceName/context/:shkeptncontext', component: ProjectBoardComponent},
  {path: 'project/:projectName/service/:serviceName/context/:shkeptncontext/stage/:stage', component: ProjectBoardComponent},
  {path: 'project/:projectName/sequence', component: ProjectBoardComponent},
  {path: 'project/:projectName/sequence/:shkeptncontext', component: ProjectBoardComponent},
  {path: 'project/:projectName/sequence/:shkeptncontext/event/:eventId', component: ProjectBoardComponent},
  {path: 'project/:projectName/sequence/:shkeptncontext/stage/:stage', component: ProjectBoardComponent},
  {path: 'trace/:shkeptncontext', component: ProjectBoardComponent},
  {path: 'trace/:shkeptncontext/:eventselector', component: ProjectBoardComponent},
  {path: 'evaluation/:shkeptncontext', component: EvaluationBoardComponent},
  {path: 'evaluation/:shkeptncontext/:eventselector', component: EvaluationBoardComponent},
  {path: 'project/:projectName/:serviceName', component: ProjectBoardComponent, canActivate: [ForwarderGuard]}, // deprecated
  {path: 'project/:projectName/:serviceName/:contextId', component: ProjectBoardComponent, canActivate: [ForwarderGuard]}, // deprecated
  {path: 'project/:projectName/:serviceName/:contextId/:eventId', component: ProjectBoardComponent, canActivate: [ForwarderGuard]}, // deprecated
  {path: '**', redirectTo: ''}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
class AppRouting {
}

export {AppRouting, routes};
