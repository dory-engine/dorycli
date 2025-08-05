#!/usr/bin/env bash
# Usage: ./pods-ready.sh <namespace>
# The default namespace is used by default when no parameters are passed

set -euo pipefail

NS="${1:-default}"

echo "Waiting for all Pods in namespace \"$NS\" to be ready"

while true; do
  pods=$(kubectl get pods -n "$NS" --no-headers 2>/dev/null || true)

  [[ -z "$pods" ]] && { echo "(No Pod)"; sleep 5; continue; }

  all_ready=true

  while IFS= read -r line; do
    ready_col=$(awk '{print $2}' <<< "$line")
    status_col=$(awk '{print $3}' <<< "$line")

    ready_cnt=${ready_col%%/*}
    total_cnt=${ready_col##*/}

    if [[ "$status_col" != "Running" || "$ready_cnt" != "$total_cnt" ]]; then
      all_ready=false
      break
    fi
  done <<< "$pods"

  kubectl get pods -n "$NS"

  if $all_ready; then
  echo "All Pods in namespace \"$NS\" are Running and Ready"
  exit 0
  fi

  echo "There are still Pods that are not ready. Check again in 5 seconds..."
  sleep 5
done
