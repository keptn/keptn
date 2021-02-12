import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'keptnUrl'
})
export class KeptnUrlPipe implements PipeTransform {

  transform(relativePath: string): string {
    return `https://keptn.sh/docs/0.8.x${relativePath}`;
  }

}
