import { moduleMetadata } from '@storybook/angular';
import { Meta, Story } from '@storybook/angular/types-6-0';
import { KtbWebhookSettingsComponent } from '../../client/app/_components/ktb-webhook-settings/ktb-webhook-settings.component';
import { KtbWebhookSettingsModule } from '../../client/app/_components/ktb-webhook-settings/ktb-webhook-settings.module';

export default {
  title: 'Components/Webhook Settings',
  component: KtbWebhookSettingsComponent,
  decorators: [
    moduleMetadata({
      imports: [KtbWebhookSettingsModule],
    }),
  ],
} as Meta;

const template: Story<KtbWebhookSettingsComponent> = (args: KtbWebhookSettingsComponent) => ({
  props: args,
  template: `<ktb-webhook-settings [webhook]="args.webhook"></ktb-webhook-settings>`,
});

export const standard = template.bind({});
standard.args = {
  webhook: {
    type: 'MYTYPE',
    method: 'GET',
    url: 'http://webhook',
    payload: '',
    header: [],
    sendFinished: true,
    sendStarted: true,
  },
};
