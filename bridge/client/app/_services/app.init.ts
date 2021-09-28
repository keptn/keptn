import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { WindowConfig } from '../../environments/environment.dynamic';

declare global {
  interface Window {
    config: WindowConfig;
  }
}

@Injectable()
export class AppInitService {
  public init(): Promise<null | string> {
    return new Promise((resolve) => {
      fetch(environment.appConfigUrl)
        .then((response) => response.text())
        .then((config) => {
          try {
            if (config) {
              Object.defineProperty(window, 'config', {
                value: JSON.parse(config),
              });

              if (window.config?.stylesheetUrl) {
                const head = document.getElementsByTagName('head')[0];
                const link = document.createElement('link');
                link.setAttribute('rel', 'stylesheet');
                link.setAttribute('type', 'text/css');
                link.setAttribute('href', window.config.stylesheetUrl);
                link.setAttribute('media', 'all');
                head.appendChild(link);
              }
            }
          } catch (err) {
            console.error('Error parsing app-config.json:', err);
          }

          return resolve(config);
        })
        .catch((err) => {
          console.error('Error loading app-config.json.', err);
          return resolve(null);
        });
    });
  }
}
