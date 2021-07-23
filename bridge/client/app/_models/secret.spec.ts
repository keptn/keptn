import { Secret } from './secret';
import { waitForAsync } from '@angular/core/testing';

describe('Secret', () => {

  it('should create with default scope "keptn-default"', waitForAsync(() => {
    // given
    const secret: Secret = new Secret();

    // then
    expect(secret).toBeTruthy();
    expect(secret.scope).toBe('keptn-default');
    expect(secret.data.length).toEqual(0, 'Secret should have no key-value pairs initially');
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
    expect(secret.data.length).toBe(2, 'Secret should have 2 key-value pairs');
    expect(secret.getData(0).key).toBe('key1', 'Key at index 0 should be "key1"');
    expect(secret.getData(0).value).toBe('value1', 'Value at index 0 should be "value1"');
    expect(secret.getData(1).key).toBe('key2', 'Key at index 1 should be "key2"');
    expect(secret.getData(1).value).toBe('value2', 'Value at index 1 should be "value2"');
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
    expect(secret.data.length).toBe(2, 'Secret should have 2 key-value pairs');
    expect(secret.getData(0).key).toBe('key1', 'Key at index 0 should be "key1"');
    expect(secret.getData(0).value).toBe('value1', 'Value at index 0 should be "value1"');
    expect(secret.getData(1).key).toBe('key3', 'Key at index 1 should be "key3"');
    expect(secret.getData(1).value).toBe('value3', 'Value at index 1 should be "value3"');
  }));
});
