import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { FormUtils } from '../../_utils/form.utils';

@Component({
  selector: 'ktb-certificate-input',
  templateUrl: './ktb-certificate-input.component.html',
  styleUrls: [],
})
export class KtbCertificateInputComponent {
  public readonly allowedExtensions = ['.pem', '.crt', '.cer', 'der'];
  public certificateControl = new FormControl('', FormUtils.isCertificateValidator);
  public certificateForm = new FormGroup({
    certificate: this.certificateControl,
  });
  public dropError?: string;

  @Input()
  public isLoading = false;

  @Input()
  set certificateInput(certificate: string | undefined) {
    this.certificateControl.setValue(atob(certificate ?? ''));
    if (certificate !== undefined) {
      this.certificateControl.markAsDirty();
    }
  }

  @Output()
  public certificateChange = new EventEmitter<string | undefined>();

  public async validateCertificate(files: FileList | null): Promise<void> {
    const file = files?.[0];
    if (file) {
      this.certificateControl.setValue((await file.text()).trim());
      this.certificateControl.markAsDirty();
      this.certificateChanged();
    }
  }

  public certificateChanged(): void {
    this.certificateChange.emit(this.certificateForm.valid ? btoa(this.certificateControl.value) : undefined);
  }
}
