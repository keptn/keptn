import {Directive, Output, EventEmitter, HostListener, HostBinding, Input} from '@angular/core';
import {FormUtils} from '../../_utils/form.utils';

@Directive({
  selector: '[ktbDragAndDrop]'
})
export class KtbDragAndDropDirective {
  private readonly BASE_STYLE_CLASS = 'ktb-drag-and-drop p-3 pb-4';

  @Input()
  multiple = false;

  @Input()
  allowedExtensions: string[] = [];

  @Output()
  dropped: EventEmitter<FileList> = new EventEmitter();

  @Output()
  dropError: EventEmitter<string> = new EventEmitter();

  @HostBinding('class')
  private styleClass = this.BASE_STYLE_CLASS;

  @HostListener('dragover', ['$event'])
  public onDragOver(evt: DragEvent) {
    evt.preventDefault();
    evt.stopPropagation();

    this.styleClass = this.BASE_STYLE_CLASS + ' drag-over';
  }

  @HostListener('dragleave', ['$event'])
  public onDragOut(evt: DragEvent) {
    evt.preventDefault();
    evt.stopPropagation();

    this.styleClass = this.BASE_STYLE_CLASS;
  }

  @HostListener('drop', ['$event'])
  public onDrop(evt: DragEvent) {
    // if (evt.preventDefault) {
    evt.preventDefault();
    evt.stopPropagation();
    // }
    const files: FileList | undefined = evt.dataTransfer?.files;
    if (files) {
      this.styleClass = this.BASE_STYLE_CLASS;

      if (!this.multiple && files.length > 1) {
        this.dropError.emit('Please select only one file');
        return;
      }

      if (!FormUtils.isFile(files[0])) {
        this.dropError.emit('Please select only files');
        return;
      }

      if (!FormUtils.isValidFileExtensions(this.allowedExtensions, files)) {
        this.dropError.emit(`Only ${this.allowedExtensions.join(', ')} files allowed`);
        return;
      }
      this.dropped.emit(files);
      this.dropError.emit('');
    }
  }
}
