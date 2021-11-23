import { Component, Input, EventEmitter, Output } from '@angular/core';
import { Secret } from '../../_models/secret';
import { SelectTreeNode, TreeListSelectOptions } from '../ktb-tree-list-select/ktb-tree-list-select.component';
import { AbstractControl } from '@angular/forms';

@Component({
  selector: 'ktb-secret-selector',
  templateUrl: './ktb-secret-selector.component.html',
})
export class KtbSecretSelectorComponent {
  @Input() public control: AbstractControl | undefined;
  @Input() public selectionStart: number | null = null;

  @Output() changed: EventEmitter<void> = new EventEmitter<void>();

  public secretDataSource: SelectTreeNode[] = [];
  public secretOptions: TreeListSelectOptions = {
    headerText: 'Select secret',
    emptyText:
      'No secrets can be found.<p>Secrets can be configured under the menu entry "Secrets" in the Uniform.</p>',
  };

  @Input()
  set secrets(secrets: Secret[] | undefined) {
    if (secrets) {
      this.secretDataSource = secrets.map((secret: Secret) => this.mapSecret(secret));
    }
  }

  public setSecret(secret: string): void {
    const secretVar = `{{.secret.${secret}}}`;
    const firstPart = this.control?.value.slice(0, this.selectionStart);
    const secondPart = this.control?.value.slice(this.selectionStart, this.control?.value.length);
    const finalString = firstPart + secretVar + secondPart;

    this.control?.setValue(finalString);
    // Input event detection is not working reliable for adding secrets, so we have to call it to work properly
    this.changed.emit();
  }

  private mapSecret(secret: Secret): SelectTreeNode {
    const scrt: SelectTreeNode = { name: secret.name };
    if (secret.keys) {
      scrt.keys = secret.keys.map((key: string) => {
        return { name: key, path: `${secret.name}.${key}` };
      });
      scrt.keys.sort((a, b) => a.name.localeCompare(b.name));
    }
    return scrt;
  }
}
