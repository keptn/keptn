function verify_test_step() {
  if [[ $1 != '0' ]]; then
    print_error "$2"
    print_error "Keptn Test failed."
    exit 1
  fi
}