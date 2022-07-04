import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component } from '@angular/core';
import { filter } from 'rxjs/operators';
import { DataService } from './_services/data.service';

// eslint-disable-next-line @typescript-eslint/no-explicit-any,@typescript-eslint/naming-convention
declare let dT_: any;

@Component({
  selector: 'ktb-app',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  private keptnInfo$ = this.dataService.keptnInfo.pipe(filter((keptnInfo) => !!keptnInfo));

  constructor(private http: HttpClient, private dataService: DataService) {
    if (typeof dT_ !== 'undefined' && dT_.initAngularNg) {
      dT_.initAngularNg(http, HttpHeaders);
    }
    this.keptnInfo$.subscribe(() => this.dataService.loadProjects());
    this.dataService.loadKeptnInfo();
  }
}
