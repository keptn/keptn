import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { filter, take } from 'rxjs/operators';
import { DataService } from './_services/data.service';

// eslint-disable-next-line @typescript-eslint/no-explicit-any,@typescript-eslint/naming-convention
declare let dT_: any;

@Component({
  selector: 'ktb-app',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  constructor(private http: HttpClient, private dataService: DataService) {
    if (typeof dT_ !== 'undefined' && dT_.initAngularNg) {
      dT_.initAngularNg(http, HttpHeaders);
    }
  }

  public ngOnInit(): void {
    this.dataService.loadKeptnInfo();
    this.dataService.keptnInfo
      .pipe(filter((keptnInfo) => !!keptnInfo))
      .pipe(take(1))
      .subscribe(() => {
        this.dataService.loadProjects();
      });
  }
}
