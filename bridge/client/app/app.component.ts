import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Component, OnInit} from '@angular/core';
import {filter, take} from 'rxjs/operators';
import {DataService} from './_services/data.service';

declare var dT_: any;

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {

  constructor(private http: HttpClient, private dataService: DataService) {
    if (typeof dT_ !== 'undefined' && dT_.initAngularNg) {
      dT_.initAngularNg(http, HttpHeaders);
    }
  }

  ngOnInit(): void {
    this.dataService.loadKeptnInfo();
    this.dataService.keptnInfo
      .pipe(filter(keptnInfo => !!keptnInfo))
      .pipe(take(1))
      .subscribe(() => {
        this.dataService.loadProjects();
      });
  }
}
