import {AbstractControl, ValidatorFn} from "@angular/forms";

export class FormUtils {
  public static projectNameExistsValidator(projectNames: string[]): ValidatorFn {
    return (control: AbstractControl): { [key: string]: any } | null => {
      const project = projectNames.includes(control.value);
      return project ? {projectName: {value: project}} : null;
    }
  }
}
