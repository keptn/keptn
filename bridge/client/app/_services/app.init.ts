import { Injectable } from '@angular/core';
import { WindowConfig } from '../../environments/environment.dynamic';
import { ApiService } from './api.service';
import { environment } from '../../environments/environment';

@Injectable()
export class AppInitService {
  constructor(private apiService: ApiService) {}

  public init(): Promise<null | WindowConfig> {
    return new Promise((resolve) => {
      this.apiService.getLookAndFeelConfig().subscribe(
        (config) => {
          if (!config) {
            return;
          }
          environment.config = config;

          if (!config.stylesheetUrl) {
            return;
          }
          const body = document.getElementsByTagName('body')[0];
          const link = document.createElement('link');
          link.setAttribute('rel', 'stylesheet');
          link.setAttribute('type', 'text/css');
          link.setAttribute('href', config.stylesheetUrl);
          link.setAttribute('media', 'all');
          body.appendChild(link);

          return resolve(config);
        },
        (err) => {
          console.error('Error loading app-config.json.', err);
          resolve(null);
        }
      );
    });
  }
}
