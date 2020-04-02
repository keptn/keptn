import {HttpClient, HttpHeaders} from "@angular/common/http";
import {Component, OnInit} from '@angular/core';
import {ApiService} from "./_services/api.service";

declare var dT_;

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {

  public version: string;
  public notes;

  constructor(private http: HttpClient, private apiService: ApiService) {
    if(typeof dT_!='undefined' && dT_.initAngularNg){dT_.initAngularNg(http, HttpHeaders);}
    this.notes = [];
  }

  ngOnInit(): void {
    this.apiService.getVersion()
      .subscribe((response: any) => {
        this.version = response.version;
        this.showNote();
      });
  }

  showNote() {
    this.notes.push({
      type: "warning",
      message: "Test"
    });
  }
}
