import { IServiceSecret, SecretKeyValuePair } from '../../../shared/interfaces/secret';

export function addData(secret: IServiceSecret, key: string, value: string): void {
  secret.data.push({ key, value });
}

export function getData(secret: IServiceSecret, index: number): SecretKeyValuePair {
  return secret.data[index];
}

export function removeData(secret: IServiceSecret, index: number): void {
  secret.data?.splice(index, 1);
}
