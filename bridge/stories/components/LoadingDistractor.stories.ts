import { Meta, Story } from '@storybook/angular/types-6-0';
import { moduleMetadata } from '@storybook/angular';
import { KtbLoadingDistractorComponent } from '../../client/app/_components/ktb-loading/ktb-loading-distractor.component';
import { KtbLoadingModule } from '../../client/app/_components/ktb-loading/ktb-loading.module';

export default {
  title: 'Components/Loading Distractor',
  component: KtbLoadingDistractorComponent,
  decorators: [
    moduleMetadata({
      imports: [KtbLoadingModule],
    }),
  ],
} as Meta;

interface LoadingDistractorArgs {
  label: string;
  className?: string;
}

const defaultArgs: LoadingDistractorArgs = {
  label: 'Loading ...',
};

const distractorTemplate: Story<LoadingDistractorArgs> = (args: LoadingDistractorArgs) => ({
  props: args,
  template: `<ktb-loading-distractor class="${args.className ?? ''}">${args.label}</ktb-loading-distractor>`,
});

const spinnerTemplate: Story<LoadingDistractorArgs> = (args: LoadingDistractorArgs) => ({
  props: args,
  template: `<ktb-loading-spinner class="${args.className ?? ''}"></ktb-loading-spinner>`,
});

export const distractor = distractorTemplate.bind({});
distractor.args = defaultArgs;

export const distractorSmaller = distractorTemplate.bind({});
distractorSmaller.args = {
  ...defaultArgs,
  className: 'smaller',
};

export const spinner = spinnerTemplate.bind({});

export const spinnerSmaller = spinnerTemplate.bind({});
spinnerSmaller.args = {
  className: 'smaller',
};
