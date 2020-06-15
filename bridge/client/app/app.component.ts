import {HttpClient, HttpHeaders} from "@angular/common/http";
import {Component, OnInit} from '@angular/core';
import {ApiService} from "./_services/api.service";
import {NotificationsService} from "./_services/notifications.service";
import {NotificationType} from "./_models/notification";

declare var dT_;

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {

  constructor(private http: HttpClient) {
    if(typeof dT_!='undefined' && dT_.initAngularNg){dT_.initAngularNg(http, HttpHeaders);}
  }

  ngOnInit(): void {
  }

}
