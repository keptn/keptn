import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';

@Component({
  selector: 'ktb-project-settings-shipyard',
  templateUrl: './ktb-project-settings-shipyard.component.html',
  styleUrls: ['./ktb-project-settings-shipyard.component.scss']
})
export class KtbProjectSettingsShipyardComponent implements OnInit {

  @Input()
  public isCreateMode: boolean;

  @Output()
  private shipyardFileChanged: EventEmitter<File> = new EventEmitter();

  public shipyardFile: File;

  constructor() { }

  ngOnInit(): void {
  }

  public updateFile(files: FileList) {
    if (files) {
      this.shipyardFile = files[0];
    } else {
      this.shipyardFile = null;
    }
    this.shipyardFileChanged.emit(this.shipyardFile);
  }
}
