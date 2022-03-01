import { Secret } from './secret';

describe('Secret', () => {
  it('should create with empty scope', () => {
    // given
    const secret: Secret = new Secret();

    // then
    expect(secret).toBeTruthy();
    expect(secret.scope).toBe('');
    expect(secret.data?.length).toEqual(0);
  });

  it('should set name property correctly', () => {
    // given
    const secret: Secret = new Secret();

    // when
    secret.setName('secretName');

    // then
    expect(secret.name).toBe('secretName');
  });

  it('should store key-value pairs', () => {
    // given
    const secret: Secret = new Secret();

    // when
    secret.addData('key1', 'value1');
    secret.addData('key2', 'value2');

    // then
    expect(secret.data?.length).toEqual(2);
    expect(secret.getData(0).key).toEqual('key1');
    expect(secret.getData(0).value).toEqual('value1');
    expect(secret.getData(1).key).toEqual('key2');
    expect(secret.getData(1).value).toEqual('value2');
  });

  it('should remove key-value pairs', () => {
    // given
    const secret: Secret = new Secret();

    // when
    secret.addData('key1', 'value1');
    secret.addData('key2', 'value2');
    secret.addData('key3', 'value3');
    secret.removeData(1);

    // then
    expect(secret.data?.length).toEqual(2);
    expect(secret.getData(0).key).toEqual('key1');
    expect(secret.getData(0).value).toEqual('value1');
    expect(secret.getData(1).key).toEqual('key3');
    expect(secret.getData(1).value).toEqual('value3');
  });
});
