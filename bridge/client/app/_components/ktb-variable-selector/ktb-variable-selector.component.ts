import { Component, EventEmitter, Input, Output, TemplateRef } from '@angular/core';
import { SelectTreeNode, TreeListSelectOptions } from '../ktb-tree-list-select/ktb-tree-list-select.component';
import { AbstractControl } from '@angular/forms';
import { DtIconType } from '@dynatrace/barista-icons';

@Component({
  selector: 'ktb-variable-selector',
  templateUrl: './ktb-variable-selector.component.html',
})
export class KtbVariableSelectorComponent {
  public options: TreeListSelectOptions = {
    headerText: 'Select element',
    emptyTemplate: null,
    hintText: '',
  };
  @Output() changed: EventEmitter<void> = new EventEmitter<void>();

  @Input() public control: AbstractControl | undefined;
  @Input() public selectionStart: number | null = null;
  @Input() public iconName: DtIconType = 'resetpassword';
  @Input() public label = '';
  @Input() public data: SelectTreeNode[] | undefined;

  @Input()
  set title(title: string) {
    this.options.headerText = title;
  }

  @Input()
  set emptyTemplate(template: TemplateRef<any>) {
    this.options.emptyTemplate = template;
  }

  @Input()
  set hintText(text: string) {
    this.options.hintText = text;
  }

  public setVariable(variable: string): void {
    if (this.control) {
      const firstPart = this.control.value.slice(0, this.selectionStart);
      const secondPart = this.control.value.slice(this.selectionStart);
      const finalString = `${firstPart}{{${variable}}}${secondPart}`;

      this.control.setValue(finalString);
      // Input event detection is not working reliable for setting the value, so we have to call it to work properly
      this.changed.emit();
    }
  }
}
