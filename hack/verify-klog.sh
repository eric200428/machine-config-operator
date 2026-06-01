#!/usr/bin/env bash

# Check for klog calls with trailing \n in non-vendor code.
# As documented at https://pkg.go.dev/k8s.io/klog#section-readme,
# klog automatically appends newlines, so explicit \n causes double newlines.

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
cd "${REPO_ROOT}"

# Look for klog calls with trailing \n" pattern
ret=0
result=$(grep -rn 'klog\.\(Info\|Infof\|Warning\|Warningf\|Error\|Errorf\|Fatal\|Fatalf\|V([0-9])\).*\\n"' \
  --include="*.go" \
  --exclude-dir=vendor \
  --exclude-dir=_output \
  . || true)

if [[ -n "${result}" ]]; then
  echo "ERROR: Found klog calls with trailing \\n (klog adds newlines automatically):"
  echo "${result}"
  echo ""
  echo "Please remove trailing \\n from the log message strings."
  echo "See https://pkg.go.dev/k8s.io/klog#section-readme for details."
  ret=1
fi

if [[ $ret -eq 0 ]]; then
  echo "Verified: No klog calls with trailing \\n found"
else
  exit 1
fi
