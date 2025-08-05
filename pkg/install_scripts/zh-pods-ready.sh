#!/usr/bin/env bash
# 用法： ./pods-ready.sh <namespace>
# 不传参时默认使用 default 命名空间

set -euo pipefail

NS="${1:-default}"

echo "正在等待 namespace \"$NS\" 中所有 Pod 就绪"

while true; do
  pods=$(kubectl get pods -n "$NS" --no-headers 2>/dev/null || true)

  [[ -z "$pods" ]] && { echo "(暂无 Pod)"; sleep 5; continue; }

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
    echo "namespace \"$NS\" 中的所有 Pod 均已 Running 且 Ready"
    exit 0
  fi

  echo "仍有 Pod 未就绪，5 秒后继续检查…"
  sleep 5
done
