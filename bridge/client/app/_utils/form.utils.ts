import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

export class FormUtils {
  public static nameExistsValidator(names: string[]): ValidatorFn {
    return (control: AbstractControl): { duplicate: { value: boolean } } | null => {
      const name = names.includes(control.value);
      return name ? { duplicate: { value: name } } : null;
    };
  }

  public static isUrlValidator(
    control: AbstractControl
  ): { url: { value: boolean } } | { space: { value: boolean } } | null {
    if (control.value) {
      if (control.value.search(/^http(s?):\/\//) === -1) {
        return { url: { value: true } };
      } else if (control.value.includes(' ')) {
        return { space: { value: true } };
      }
    }
    return null;
  }

  public static isUrlOrSecretValidator(
    control: AbstractControl
  ): { urlOrSecret: { value: boolean } } | { space: { value: boolean } } | null {
    if (control.value) {
      if ((control.value.search(/^http(s?):\/\//) === -1) && (control.value.search(/^{{(.+)}}/) === -1)) {
        return { urlOrSecret: { value: true } };
      } else if (control.value.includes(' ')) {
        return { space: { value: true } };
      }
    }
    return null;
  }

  public static isSshValidator(control: AbstractControl): ValidationErrors | null {
    if (control.value && !control.value.startsWith('ssh://')) {
      return { ssh: true };
    }
    return null;
  }

  public static payloadSpecialCharValidator(control: AbstractControl): { specialChar: { value: boolean } } | null {
    if (control.value && control.value.match(/(\$|\||;|>|&|`|\/var\/run)/gi)) {
      return { specialChar: { value: true } };
    }
    return null;
  }

  public static isCertificateValidator(control: AbstractControl): ValidationErrors | null {
    if (
      control.value &&
      (!control.value.startsWith('-----BEGIN CERTIFICATE-----') || !control.value.endsWith('-----END CERTIFICATE-----'))
    ) {
      return { certificate: true };
    }
    return null;
  }

  public static isSshKeyValidator(control: AbstractControl): ValidationErrors | null {
    if (
      control.value &&
      (!control.value.startsWith('-----BEGIN OPENSSH PRIVATE KEY-----') ||
        !control.value.endsWith('-----END OPENSSH PRIVATE KEY-----'))
    ) {
      return { sshKey: true };
    }
    return null;
  }

  public static isValidFileExtensions(allowedExtensions: string[], files: FileList): boolean {
    if (allowedExtensions && allowedExtensions.length > 0) {
      const allowedFiles = [];
      allowedExtensions.forEach((extension) => {
        const fileArray: File[] = Array.from(files);
        fileArray.forEach((file) => {
          if (file.name.endsWith(extension)) {
            allowedFiles.push(file);
          }
        });
      });
      if (allowedFiles.length === 0) {
        return false;
      }
    }
    return true;
  }

  public static isFile(file: File): boolean {
    return !(!file.type && file.size % 4096 === 0);
  }

  public static readFileContent(file: File): Promise<string | undefined> {
    return new Promise<string | undefined>((resolve, reject) => {
      if (!file) {
        resolve('');
      }

      const reader = new FileReader();
      reader.onload = (): void => {
        try {
          const text = reader.result?.toString();
          resolve(text);
        } catch (e) {
          reject(e);
        }
      };

      reader.readAsText(file);
    });
  }
}
