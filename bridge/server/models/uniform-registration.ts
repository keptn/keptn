import { UniformRegistration as ur } from '../../shared/models/uniform-registration';

export class UniformRegistration extends ur {
  public static fromJSON(data: unknown): UniformRegistration {
    return Object.assign(new this(), data);
  }
}
