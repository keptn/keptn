import { Secret } from './secret';
import { waitForAsync } from '@angular/core/testing';

describe('Secret', () => {

  it('should create with default scope "keptn-default"', waitForAsync(() => {
    // given
    const secret: Secret = new Secret();

    // then
    expect(secret).toBeTruthy();
    expect(secret.scope).toBe('keptn-default');
    expect(secret.data.length).toEqual(0);
  }));

  it('should set name property correctly', waitForAsync(() => {
    // given
    const secret: Secret = new Secret();

    // when
    secret.setName('secretName');

    // then
    expect(secret.name).toBe('secretName');
  }));

  it('should store key-value pairs', waitForAsync(() => {
    // given
    const secret: Secret = new Secret();

    // when
    secret.addData('key1', 'value1');
    secret.addData('key2', 'value2');

    // then
    expect(secret.data.length).toEqual(2);
    expect(secret.getData(0).key).toEqual('key1');
    expect(secret.getData(0).value).toEqual('value1');
    expect(secret.getData(1).key).toEqual('key2');
    expect(secret.getData(1).value).toEqual('value2');
  }));

  it('should remove key-value pairs', waitForAsync(() => {
    // given
    const secret: Secret = new Secret();

    // when
    secret.addData('key1', 'value1');
    secret.addData('key2', 'value2');
    secret.addData('key3', 'value3');
    secret.removeData(1);

    // then
    expect(secret.data.length).toEqual(2);
    expect(secret.getData(0).key).toEqual('key1');
    expect(secret.getData(0).value).toEqual('value1');
    expect(secret.getData(1).key).toEqual('key3');
    expect(secret.getData(1).value).toEqual('value3');
  }));
});
