import { Meta, Story } from '@storybook/angular/types-6-0';
import { moduleMetadata } from '@storybook/angular';
import { KtbDatetimePickerDirective } from '../../client/app/_components/ktb-date-input/ktb-datetime-picker.component';
import { KtbDateInputModule } from '../../client/app/_components/ktb-date-input/ktb-date-input.module';

import '@angular/localize/init';
import { EventEmitter } from '@angular/core';
import { RouterTestingModule } from '@angular/router/testing';
import { DtButtonModule } from '@dynatrace/barista-components/button';

export default {
  title: 'Components/Datetime Picker',
  component: KtbDatetimePickerDirective,
  decorators: [
    moduleMetadata({
      imports: [KtbDateInputModule, DtButtonModule, RouterTestingModule],
    }),
  ],
} as Meta;

const template: Story = (args) => ({
  props: args,
  template: `<button dt-button ktbDatetimePicker>Pick date/time</button>`,
});

export const standard = template.bind({});
standard.args = {
  timeEnabled: true,
  secondsEnabled: false,
  closeDialog: new EventEmitter<void>(),
};
