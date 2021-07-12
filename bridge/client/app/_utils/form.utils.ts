import {AbstractControl, ValidatorFn} from "@angular/forms";

export class FormUtils {
  public static projectNameExistsValidator(projectNames: string[]): ValidatorFn {
    return (control: AbstractControl): { [key: string]: any } | null => {
      const project = projectNames.includes(control.value);
      return project ? {projectName: {value: project}} : null;
    }
  }

  public static isValidFileExtensions(allowedExtensions: string[], files: FileList): boolean {
    if (allowedExtensions && allowedExtensions.length > 0) {
      const allowedFiles = [];
      allowedExtensions.forEach(extension => {
        const fileArray: File[] = Array.from(files);
        fileArray.forEach(file => {
          if(file.name.endsWith(extension)) {
            allowedFiles.push(file);
          }
        });
      });
      if(allowedFiles.length === 0) {
        return false;
      }
    }
    return true;
  }

  public static isFile(file: File): boolean {
    return !(!file.type && file.size % 4096 == 0);
  }

  public static readFileContent(file: File): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      if (!file) {
        resolve('');
      }

      const reader = new FileReader();
      reader.onload = () => {
        try {
          const text = reader.result.toString();
          resolve(text);
        } catch (e) {
          reject(e);
        }
      };

      reader.readAsText(file);
    });
  }
}
