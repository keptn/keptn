export default class SearchUtil {
  static binarySearch(array, element, compare_fn) {
    if(compare_fn(element, array[0]) < 0) // compare with first element and insert at beginning if lower
      return 0;
    if(compare_fn(element, array[array.length-1]) > 0) // compare with last element and insert at the end if greater
      return array.length;

    let lower_index = 0;
    let upper_index = array.length - 1;
    while (lower_index <= upper_index) { // find position to insert with binarySearch logic
      let median_index = (upper_index + lower_index) >> 1; // determine median index
      let cmp_result = compare_fn(element, array[median_index]); // compare element with median element

      if (cmp_result > 0) { // if greater then median element, compare with second half in next iteration
        lower_index = median_index + 1;
      } else if(cmp_result < 0) { // if smaller then median element, compare with first half in next iteration
        upper_index = median_index - 1;
      } else { // if equal with median, insert at same index
        return median_index;
      }
    }
    return lower_index; // lower_index passed upper_index and is therefore the position to insert element
  }
}
