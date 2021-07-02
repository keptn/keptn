import {Directive, Output, EventEmitter, HostListener, HostBinding, Input} from '@angular/core';

@Directive({
  selector: '[ktbDragAndDrop]'
})
export class KtbDragAndDropDirective {
  private readonly BASE_STYLE_CLASS = 'ktb-drag-and-drop p-3 pb-4';

  @Input()
  multiple: boolean = false;

  @Input()
  allowedExtensions: string[];

  @Output()
  onDropped: EventEmitter<FileList> = new EventEmitter();

  @Output()
  onError: EventEmitter<string> = new EventEmitter();

  @HostBinding('class')
  private styleClass = this.BASE_STYLE_CLASS;

  @HostListener('dragover', ['$event'])
  public onDragOver(evt) {
    evt.preventDefault();
    evt.stopPropagation();

    this.styleClass = this.BASE_STYLE_CLASS + ' drag-over';
  }

  @HostListener('dragleave', ['$event'])
  public onDragOut(evt) {
    evt.preventDefault();
    evt.stopPropagation();

    this.styleClass = this.BASE_STYLE_CLASS;
  }

  @HostListener('drop', ['$event'])
  public onDrop(evt) {
    evt.preventDefault();
    evt.stopPropagation();
    const files: FileList = evt.dataTransfer.files;
    this.styleClass = this.BASE_STYLE_CLASS;

    if(!this.multiple && files.length > 1) {
      this.onError.emit('Please select only one file');
      return;
    }

    if (!files[0].type && files[0].size%4096 == 0) {
      this.onError.emit('Please select only files');
      return;
    }

    if (this.allowedExtensions && this.allowedExtensions.length > 0) {
      const allowedFiles = [];
      this.allowedExtensions.forEach(extension => {
        const fileArray: File[] = Array.from(files);
        fileArray.forEach(file => {
          if(file.name.endsWith(extension)) {
            allowedFiles.push(file);
          }
        });
      });
      if(allowedFiles.length === 0) {
        this.onError.emit(`Only ${this.allowedExtensions.join(', ')} files allowed`);
        return;
      }
    }

    this.onDropped.emit(files);
    this.onError.emit('');
  }

  constructor() {}
}
