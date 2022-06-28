import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { KeyValue } from '@angular/common';
import { AppUtils } from '../../_utils/app.utils';
import { IProxy } from '../../../../shared/interfaces/Project';

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
  public schemeControl = new FormControl(this.schemes[1].value);
  public passwordControl = new FormControl('');
  public userControl = new FormControl('');
  public hostControl = new FormControl('', [Validators.required]);
  public portControl = new FormControl('', [Validators.required]);
  public proxyForm = new FormGroup({
    scheme: this.schemeControl,
    user: this.userControl,
    password: this.passwordControl,
    host: this.hostControl,
    port: this.portControl,
  });
  @Input()
  public set proxy(proxy: IProxy | undefined) {
    if (proxy) {
      const urlParts = AppUtils.splitURLPort(proxy.url);
      this.schemeControl.setValue(proxy.scheme);
      this.hostControl.setValue(urlParts.host);
      this.portControl.setValue(urlParts.port);
      this.userControl.setValue(proxy.user ?? '');
      this.passwordControl.setValue(proxy.password ?? '');
    }
  }
  public get proxy(): IProxy | undefined {
    return this.proxyForm.valid
      ? {
          scheme: this.schemeControl.value,
          url: `${this.hostControl.value}:${this.portControl.value}`,
          user: this.userControl.value,
          password: this.passwordControl.value,
        }
      : undefined;
  }
  @Output()
  public proxyChange = new EventEmitter<IProxy | undefined>();

  public proxyChanged(): void {
    this.proxyChange.emit(this.proxy);
  }
}
