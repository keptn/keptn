import { Secret } from './secret';
import { waitForAsync } from '@angular/core/testing';

describe('Secret', () => {

  it('should create with default scope "keptn-default"', waitForAsync(() => {
    const secret: Secret = new Secret();

    expect(secret).toBeTruthy();
    expect(secret.scope).toBe('keptn-default');
  }));

  it('should store key-value pairs', waitForAsync(() => {
    const secret: Secret = new Secret();

    secret.addData('key1', 'value1');
    secret.addData('key2', 'value2');

    expect(secret.data.length).toBe(2, 'Secret should have 2 key-value pairs');
    expect(secret.getData(0).key).toBe('key1', 'Key at index 0 should be "key1"');
    expect(secret.getData(0).value).toBe('value1', 'Value at index 0 should be "value1"');
    expect(secret.getData(1).key).toBe('key2', 'Key at index 1 should be "key2"');
    expect(secret.getData(1).value).toBe('value2', 'Value at index 1 should be "value2"');
  }));
});
