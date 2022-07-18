import { Meta, Story } from '@storybook/angular/types-6-0';
import { KtbHeatmapComponent } from '../../client/app/_components/ktb-heatmap/ktb-heatmap.component';
import { generateTestData } from '../../client/app/_components/ktb-heatmap/testing/ktb-test-heatmap.component';
import { componentWrapperDecorator, moduleMetadata } from '@storybook/angular';
import { KtbHeatmapModule } from '../../client/app/_components/ktb-heatmap/ktb-heatmap.module';

export default {
  title: 'Components/Heatmap',
  component: KtbHeatmapComponent,
  decorators: [
    moduleMetadata({
      imports: [KtbHeatmapModule],
    }),
    componentWrapperDecorator((story) => `<div style="margin: 16px">${story}</div>`),
  ],
  parameters: {
    layout: 'fullscreen',
  },
} as Meta;

const template: Story<KtbHeatmapComponent> = (args: KtbHeatmapComponent) => ({
  props: args,
});

export const random = template.bind({});
random.args = {
  dataPoints: generateTestData(5, 8),
};

export const randomLarge = template.bind({});
randomLarge.args = {
  dataPoints: generateTestData(15, 40),
  showMoreVisible: true,
};
