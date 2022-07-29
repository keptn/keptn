import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { moduleMetadata } from '@storybook/angular';
import { Meta, Story } from '@storybook/angular/types-6-0';
import { KtbCertificateInputComponent } from '../../client/app/_components/ktb-certificate-input/ktb-certificate-input.component';
import { KtbCertificateInputModule } from '../../client/app/_components/ktb-certificate-input/ktb-certificate-input.module';

export default {
  title: 'Components/Certificate Input',
  decorators: [
    moduleMetadata({
      imports: [KtbCertificateInputModule, BrowserAnimationsModule],
    }),
  ],
} as Meta;

const template: Story<KtbCertificateInputComponent> = (args: KtbCertificateInputComponent) => ({
  props: args,
  template: `<ktb-certificate-input [certificateInput]="certificateInput"></ktb-certificate-input>`,
});

export const selected = template.bind({});
selected.args = {
  certificateInput: btoa('---BEGIN--- CERTIFICATE ---END---'),
};
