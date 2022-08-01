import { HttpClientModule } from '@angular/common/http';
import { moduleMetadata } from '@storybook/angular';
import { Meta, Story } from '@storybook/angular/types-6-0';
import { KtbCopyToClipboardComponent } from '../../client/app/_components/ktb-copy-to-clipboard/ktb-copy-to-clipboard.component';
import { KtbCopyToClipboardModule } from '../../client/app/_components/ktb-copy-to-clipboard/ktb-copy-to-clipboard.module';

export default {
  title: 'Components/Copy to Clipboard',
  component: KtbCopyToClipboardComponent,
  decorators: [
    moduleMetadata({
      imports: [KtbCopyToClipboardModule, HttpClientModule],
    }),
  ],
} as Meta;

const template: Story<KtbCopyToClipboardComponent> = (args: KtbCopyToClipboardComponent) => ({
  props: args,
});

export const standard = template.bind({});
standard.args = {
  value: 'My Value',
  label: 'This is my value!',
};
