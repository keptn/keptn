import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class ForwarderGuard implements CanActivate {
  constructor(private router: Router) {}

  canActivate(route: ActivatedRouteSnapshot) {
    // => project/:projectName/service/:serviceName
    // => project/:projectName/sequence/:shkeptncontext/event/:eventId
    if (route.params.eventId) {
      // project/:projectName/:serviceName/:contextId/:eventId
      this.router.navigate([
        'project',
        route.params.projectName,
        'sequence',
        route.params.contextId,
        'event',
        route.params.eventId,
      ]);
    } else if (route.params.contextId) {
      // project/:projectName/:serviceName/:contextId
      this.router.navigate(['project', route.params.projectName, 'sequence', route.params.contextId]);
    } else {
      // project/:projectName/:serviceName
      this.router.navigate(['project', route.params.projectName, 'service', route.params.serviceName]);
    }
    return false;
  }
}
