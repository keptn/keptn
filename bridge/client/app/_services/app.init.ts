import { Injectable } from '@angular/core';
import { from } from 'rxjs';
declare var window: any;

@Injectable()
export class AppInitService {

  public init() {
    return from(
      fetch('assets/branding/app-config.json', {
        headers : {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }).then(response => {
        return response.json();
      }).then(config => {
        if(config)
          window.config = config;
        return config;
      }).catch(err => {
        console.log("Error loading app-config.json " + err);
        return null;
      })
    ).toPromise();
  }
}
