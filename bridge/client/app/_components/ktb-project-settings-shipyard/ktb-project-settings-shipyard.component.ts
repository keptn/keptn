import { Component, ElementRef, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { FormUtils } from '../../_utils/form.utils';

@Component({
  selector: 'ktb-project-settings-shipyard',
  templateUrl: './ktb-project-settings-shipyard.component.html',
  styleUrls: ['./ktb-project-settings-shipyard.component.scss'],
})
export class KtbProjectSettingsShipyardComponent {
  public readonly allowedExtensions = ['yaml', 'yml'];

  @Input()
  public isCreateMode = false;

  @Output()
  public shipyardFileChanged: EventEmitter<File | undefined> = new EventEmitter();

  @ViewChild('dropError')
  private dropError?: ElementRef;

  public shipyardFile?: File;

  public handleDragAndDropError(error: string) {
    if (this.dropError) {
      this.dropError.nativeElement.innerText = error;
    }
  }

  public updateFile(files?: FileList) {
    if (files) {
      this.shipyardFile = files[0];
    } else {
      this.shipyardFile = undefined;
    }
    this.shipyardFileChanged.emit(this.shipyardFile);
  }

  public validateAndUpdateFile(files: FileList | null) {
    if (files?.length && this.dropError) {
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
