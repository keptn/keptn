import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  HostBinding,
  Input,
  ViewEncapsulation,
} from '@angular/core';
import { Platform } from '@angular/cdk/platform';

@Component({
  selector: 'ktb-copy-to-clipboard[value][label]',
  templateUrl: './ktb-copy-to-clipboard.component.html',
  styleUrls: ['./ktb-copy-to-clipboard.component.scss'],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbCopyToClipboardComponent {
  public visible = false;
  @HostBinding('class') cls = 'ktb-copy-to-clipboard';
  @HostBinding('attr.aria-visible') ariaVisible = this.visible;
  @HostBinding('class.ktb-copy-input-visible') inputVisible = this.visible;
  @Input() public value = '';
  @Input() public label = '';

  constructor(private _changeDetectorRef: ChangeDetectorRef, public platform: Platform) {}
}
