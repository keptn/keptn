import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component } from '@angular/core';
import { filter, take } from 'rxjs/operators';
import { DataService } from './_services/data.service';
import { ActivatedRoute } from '@angular/router';
import { combineLatest, of } from 'rxjs';

// eslint-disable-next-line @typescript-eslint/no-explicit-any,@typescript-eslint/naming-convention
declare let dT_: any;

@Component({
  selector: 'ktb-app',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  constructor(private http: HttpClient, private dataService: DataService, private route: ActivatedRoute) {
    if (typeof dT_ !== 'undefined' && dT_.initAngularNg) {
      dT_.initAngularNg(http, HttpHeaders);
    }
    this.dataService.loadKeptnInfo();
    const keptnInfo$ = this.dataService.keptnInfo.pipe(filter((keptnInfo) => !!keptnInfo));

    combineLatest([this.route.firstChild?.data ?? of({}), keptnInfo$])
      .pipe(take(1))
      .subscribe(([data]) => {
        if (!data.projectsHandledByComponent) {
          this.dataService.loadProjects().subscribe();
        }
      });
  }
}
