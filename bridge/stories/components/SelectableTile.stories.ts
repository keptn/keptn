import { moduleMetadata } from '@storybook/angular';
import { Meta, Story } from '@storybook/angular/types-6-0';
import { KtbSelectableTileComponent } from '../../client/app/_components/ktb-selectable-tile/ktb-selectable-tile.component';
import { KtbSelectableTileModule } from '../../client/app/_components/ktb-selectable-tile/ktb-selectable-tile.module';

export default {
  title: 'Components/Selectable Tile',
  component: KtbSelectableTileComponent,
  decorators: [
    moduleMetadata({
      imports: [KtbSelectableTileModule],
    }),
  ],
} as Meta;

const template: Story<KtbSelectableTileComponent> = (args: KtbSelectableTileComponent) => ({
  props: args,
  template: `<ktb-selectable-tile [selected]="${args.selected}">
        <span ktb-selectable-tile-header>Selectable tile header</span>
        This is the content of the tile
    </ktb-selectable-tile>`,
});

export const selected = template.bind({});
selected.args = {
  selected: true,
};
