import { parseCurl } from './curl.utils';

describe('Test curl-parser', () => {
  it('should parse curl correctly', () => {
    expect(
      parseCurl(
        `curl http://keptn.sh/asdf asdf --request GET --header 'content-type: application/json' --header 'Authorization: Bearer myToken' --proxy http://keptn.sh/proxy --data '{"data": "myData"}'`
      )
    ).toEqual({
      _: ['http://keptn.sh/asdf', 'asdf'],
      request: ['GET'],
      header: ['content-type: application/json', 'Authorization: Bearer myToken'],
      proxy: ['http://keptn.sh/proxy'],
      data: ['{"data": "myData"}'],
    });
  });

  it('should parse options correctly', () => {
    expect(parseCurl(`--proxy http://keptn.sh/proxy --data '{"data": "myData"}'`, true)).toEqual({
      proxy: ['http://keptn.sh/proxy'],
      data: ['{"data": "myData"}'],
    });
  });
});
