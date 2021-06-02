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
        return response.text();
      }).then(config => {
        try {
          if(config)
            window.config = JSON.parse(config);
          return config;
        } catch(err) {
          console.error("Error parsing app-config.json:", err);
          return null;
        }
      }).catch(err => {
        console.error("Error loading app-config.json.", err);
        return null;
      })
    ).toPromise();
  }
}
