import { AbstractControl, ValidatorFn } from '@angular/forms';

export class FormUtils {
  public static nameExistsValidator(names: string[]): ValidatorFn {
    return (control: AbstractControl): { duplicate: { value: boolean } } | null => {
      const name = names.includes(control.value);
      return name ? {duplicate: {value: name}} : null;
    };
  }

  public static isValidFileExtensions(allowedExtensions: string[], files: FileList): boolean {
    if (allowedExtensions && allowedExtensions.length > 0) {
      const allowedFiles = [];
      allowedExtensions.forEach(extension => {
        const fileArray: File[] = Array.from(files);
        fileArray.forEach(file => {
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
      reader.onload = () => {
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
