import { Component, Input, EventEmitter, Output } from '@angular/core';
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

  @Input() public data: SelectTreeNode[] = [];
  public options: TreeListSelectOptions = {
    headerText: 'Select element',
    emptyText: 'No elements available.',
  };

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
}
