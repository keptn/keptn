import {Component, ElementRef, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {FormUtils} from '../../_utils/form.utils';

@Component({
  selector: 'ktb-project-settings-shipyard',
  templateUrl: './ktb-project-settings-shipyard.component.html',
  styleUrls: ['./ktb-project-settings-shipyard.component.scss']
})
export class KtbProjectSettingsShipyardComponent implements OnInit {
  public readonly allowedExtensions = ['yaml', 'yml'];

  @Input()
  public isCreateMode: boolean;

  @Output()
  public shipyardFileChanged: EventEmitter<File> = new EventEmitter();

  @ViewChild('dropError')
  private dropError: ElementRef;

  public shipyardFile: File | null;

  constructor() {
  }

  ngOnInit(): void {
  }

  public handleDragAndDropError(error: string) {
    this.dropError.nativeElement.innerText = error;
  }

  public updateFile(files: FileList) {
    if (files) {
      this.shipyardFile = files[0];
    } else {
      this.shipyardFile = null;
    }
    this.shipyardFileChanged.emit(this.shipyardFile);
  }

  public validateAndUpdateFile(files: FileList) {
    if (files && files.length > 0) {
      if (!FormUtils.isFile(files[0])) {
        this.dropError.nativeElement.innerText = 'Please select only files';
        return;
      }

      if (!FormUtils.isValidFileExtensions(this.allowedExtensions, files)) {
        this.dropError.nativeElement.innerText = `Only ${this.allowedExtensions.join(', ')} files allowed`;
        return;
      }

      this.dropError.nativeElement.innerText = '';
      this.updateFile(files);
    }
  }
}
