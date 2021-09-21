import {Component, Input} from '@angular/core';
import {Sequence} from '../../_models/sequence';

@Component({
  selector: 'ktb-sequence-state-info',
  templateUrl: './ktb-sequence-state-info.component.html',
  styleUrls: ['./ktb-sequence-state-info.component.scss']
})
export class KtbSequenceStateInfoComponent {

  private _sequence?: Sequence;

  @Input()
  get sequence(): Sequence | undefined {
    return this._sequence;
  }
  set sequence(sequence: Sequence | undefined) {
    if (this._sequence !== sequence) {
      this._sequence = sequence;
    }
  }

  constructor() {
  }
}
