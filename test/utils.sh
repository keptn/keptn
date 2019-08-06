function timestamp() {
  date +"[%Y-%m-%d %H:%M:%S]"
}

function print_error() {
  echo "[keptn|ERROR] $(timestamp) $1"
}

function verify_test_step() {
  if [[ $1 != '0' ]]; then
    print_error "$2"
    print_error "Keptn Test failed."
    exit 1
  fi
}