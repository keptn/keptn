import { Meta, Story } from '@storybook/angular/types-6-0';
import { moduleMetadata } from '@storybook/angular';

import '@angular/localize/init';
import { KtbAppHeaderComponent } from '../../client/app/_components/ktb-app-header/ktb-app-header.component';
import { KtbAppHeaderModule } from '../../client/app/_components/ktb-app-header/ktb-app-header.module';
import { RouterTestingModule } from '@angular/router/testing';

export default {
  title: 'Components/App Header',
  component: KtbAppHeaderComponent,
  decorators: [
    moduleMetadata({
      imports: [KtbAppHeaderModule, RouterTestingModule],
    }),
  ],
} as Meta;

const template: Story<KtbAppHeaderComponent> = (args: KtbAppHeaderComponent) => ({
  props: args,
});

export const standard = template.bind({});
standard.args = {};
