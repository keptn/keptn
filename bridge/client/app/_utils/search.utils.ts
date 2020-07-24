export default class SearchUtil {
  static binarySearch(ar, el, compare_fn) {
    if(compare_fn(el, ar[0]) < 0)
      return 0;
    if(compare_fn(el, ar[ar.length-1]) > 0)
      return ar.length;
    let m = 0;
    let n = ar.length - 1;
    while (m <= n) {
      let k = (n + m) >> 1;
      let cmp = compare_fn(el, ar[k]);
      if (cmp > 0) {
        m = k + 1;
      } else if(cmp < 0) {
        n = k - 1;
      } else {
        return k;
      }
    }
    return -m - 1;
  }
}
