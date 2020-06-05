import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {DashboardComponent} from "./dashboard/dashboard.component";
import {ProjectBoardComponent} from "./project-board/project-board.component";


const routes: Routes = [
  {path: '', pathMatch: 'full', redirectTo: 'dashboard'},
  {path: 'dashboard', component: DashboardComponent},
  {path: 'project/:projectName', component: ProjectBoardComponent},
  {path: 'project/:projectName/:serviceName', component: ProjectBoardComponent},
  {path: 'project/:projectName/:serviceName/:contextId', component: ProjectBoardComponent},
  {path: 'project/:projectName/:serviceName/:contextId/:eventId', component: ProjectBoardComponent},
  {path: 'trace/:shkeptncontext', component: ProjectBoardComponent},
  {path: 'trace/:shkeptncontext/:eventselector', component: ProjectBoardComponent},
  {path: '**', redirectTo: ''}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
class AppRouting {
}

export {AppRouting, routes};
