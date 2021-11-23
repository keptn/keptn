import { Component, Input, EventEmitter, Output } from '@angular/core';
import { Secret } from '../../_models/secret';
import { SelectTreeNode, TreeListSelectOptions } from '../ktb-tree-list-select/ktb-tree-list-select.component';
import { AbstractControl } from '@angular/forms';

@Component({
  selector: 'ktb-variable-selector',
  templateUrl: './ktb-variable-selector.component.html',
})
export class KtbVariableSelectorComponent {
  @Input() public control: AbstractControl | undefined;
  @Input() public selectionStart: number | null = null;
  @Input() public variablePrefix = '';

  @Output() changed: EventEmitter<void> = new EventEmitter<void>();

  public treeDataSource: SelectTreeNode[] = [];
  public treeOptions: TreeListSelectOptions = {
    headerText: 'Select element',
    emptyText: 'No elements available.',
  };

  @Input()
  set secrets(secrets: Secret[] | undefined) {
    if (secrets) {
      this.treeDataSource = []; // secrets.map((secret: Secret) => this.mapSecret(secret));
    }
  }

  @Input()
  set title(title: string) {
    this.treeOptions.headerText = title;
  }

  @Input()
  set emptyText(text: string) {
    this.treeOptions.emptyText = text;
  }

  public setVariable(variable: string): void {
    const variableString = `{{${this.variablePrefix}.${variable}}}`;
    const firstPart = this.control?.value.slice(0, this.selectionStart);
    const secondPart = this.control?.value.slice(this.selectionStart);
    const finalString = firstPart + variableString + secondPart;

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
