import { IServiceSecret } from '../../../shared/interfaces/secret';
import { SecretScopeDefault } from '../../../shared/interfaces/secret-scope';
import { addData, getData, removeData } from './secret';

describe('Secret', () => {
  it('should return undefined data if index out of bounds', () => {
    // given
    const secret: IServiceSecret = {
      name: '',
      scope: SecretScopeDefault.DEFAULT,
      data: [],
    };

    // when
    const data = getData(secret, 999);

    // then
    expect(data).toBe(undefined);
  });

  it('should store key-value pairs', () => {
    // given
    const secret: IServiceSecret = {
      name: '',
      scope: SecretScopeDefault.DEFAULT,
      data: [],
    };

    // when
    addData(secret, 'key1', 'value1');
    addData(secret, 'key2', 'value2');

    // then
    expect(secret.data?.length).toEqual(2);
    expect(getData(secret, 0).key).toEqual('key1');
    expect(getData(secret, 0).value).toEqual('value1');
    expect(getData(secret, 1).key).toEqual('key2');
    expect(getData(secret, 1).value).toEqual('value2');
  });

  it('should remove key-value pairs', () => {
    // given
    const secret: IServiceSecret = {
      name: '',
      scope: SecretScopeDefault.DEFAULT,
      data: [],
    };

    // when
    addData(secret, 'key1', 'value1');
    addData(secret, 'key2', 'value2');
    addData(secret, 'key3', 'value3');
    removeData(secret, 1);

    // then
    expect(secret.data?.length).toEqual(2);
    expect(getData(secret, 0).key).toEqual('key1');
    expect(getData(secret, 0).value).toEqual('value1');
    expect(getData(secret, 1).key).toEqual('key3');
    expect(getData(secret, 1).value).toEqual('value3');
  });
});
