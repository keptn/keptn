export type Dimension = { width: number; height: number };
export type Position = { top: number; left: number };

const colors = [
  '#9355b7',
  '#7dc540',
  '#14a8f5',
  '#f5d30f',
  '#ef651f',
  '#dc172a',
  '#00b9cc',
  '#522273',
  '#1f7e1e',
  '#004999',
  '#ab8300',
  '#8d380f',
  '#93060e',
  '#006d75',
];

export function getColor(index: number): string {
  return colors[index % colors.length];
}

export function getIconStyle(index: number, invisible?: boolean): string {
  const color = invisible === true ? '#cccccc' : getColor(index);
  return `--dt-icon-color: ${color}`;
}

export function replaceSpace(value: string): string {
  return value.replace(/ /g, '-');
}

export function getTooltipPosition(
  windowDimensions: Dimension,
  tooltip: Dimension,
  barArea: Dimension & Position
): Position {
  let top = barArea.top + barArea.height + 20;
  let left = barArea.left;
  if (top + tooltip.height > windowDimensions.height) {
    top = windowDimensions.height - tooltip.height - 10;
  }
  if (left + tooltip.width > windowDimensions.width) {
    left = windowDimensions.width - tooltip.width - 10;
  }
  return { top, left };
}
