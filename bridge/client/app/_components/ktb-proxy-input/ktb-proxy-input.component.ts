import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { KeyValue } from '@angular/common';
import { IProxy } from '../../_interfaces/git-upstream';
import { AppUtils } from '../../_utils/app.utils';

@Component({
  selector: 'ktb-proxy-input',
  templateUrl: './ktb-proxy-input.component.html',
  styleUrls: [],
})
export class KtbProxyInputComponent {
  public readonly schemes: KeyValue<string, string>[] = [
    {
      key: 'HTTP',
      value: 'http',
    },
    {
      key: 'HTTPS',
      value: 'https',
    },
  ];
  public isInsecureControl = new FormControl(false);
  public schemeControl = new FormControl(this.schemes[1].value);
  public passwordControl = new FormControl('');
  public userControl = new FormControl('');
  public hostControl = new FormControl('', [Validators.required]);
  public portControl = new FormControl('', [Validators.required]);
  public proxyForm = new FormGroup({
    isInsecure: this.isInsecureControl,
    scheme: this.schemeControl,
    user: this.userControl,
    password: this.passwordControl,
    host: this.hostControl,
    port: this.portControl,
  });
  @Input()
  public set proxy(proxy: IProxy | undefined) {
    if (proxy) {
      const urlParts = AppUtils.splitURLPort(proxy.gitProxyUrl);
      this.isInsecureControl.setValue(proxy.gitProxyInsecure);
      this.schemeControl.setValue(proxy.gitProxyScheme);
      this.hostControl.setValue(urlParts.host);
      this.portControl.setValue(urlParts.port);
      this.userControl.setValue(proxy.gitProxyUser ?? '');
      this.passwordControl.setValue(proxy.gitProxyPassword ?? '');
    }
  }
  public get proxy(): IProxy | undefined {
    return this.proxyForm.valid
      ? {
          gitProxyInsecure: this.isInsecureControl.value,
          gitProxyScheme: this.schemeControl.value,
          gitProxyUrl: `${this.hostControl.value}:${this.portControl.value}`,
          gitProxyUser: this.userControl.value,
          gitProxyPassword: this.passwordControl.value,
        }
      : undefined;
  }
  @Output()
  public proxyChange = new EventEmitter<IProxy | undefined>();

  public proxyChanged(): void {
    this.proxyChange.emit(this.proxy);
  }
}
