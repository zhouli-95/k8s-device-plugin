#!/bin/bash

BUILD_IMAGE="${BUILD_IMAGE:-false}"
RUN_ENV="${RUN_ENV:-kind}"
CREATE_KIND_CLUSTER="${CREATE_KIND_CLUSTER:-false}"
TEST_CLUSTER_NAME="${TEST_CLUSTER_NAME:-test-cluster}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
CLUSTER_CONFIG="$PROJECT_DIR/test/config/cluster.yaml"
SAVED_KUBECONFIG=""

if [ "$BUILD_IMAGE" == "true" ]; then
  make docker
fi

if [ "$CREATE_KIND_CLUSTER" == "true" ]; then
  echo "Creating kind cluster: $TEST_CLUSTER_NAME (config: $CLUSTER_CONFIG)"
  kind create cluster --name "$TEST_CLUSTER_NAME" --config "$CLUSTER_CONFIG" --wait 2m

  SAVED_KUBECONFIG="/tmp/${TEST_CLUSTER_NAME}-kubeconfig"
  kind get kubeconfig --name "$TEST_CLUSTER_NAME" > "$SAVED_KUBECONFIG"
  chmod 600 "$SAVED_KUBECONFIG"
  echo "Kubeconfig saved to: $SAVED_KUBECONFIG"
  export KUBECONFIG="$SAVED_KUBECONFIG"
  export KIND_CLUSTER_NAME="$TEST_CLUSTER_NAME"
  bash "$SCRIPT_DIR/kind-load-image.sh"
else
  if [ "$RUN_ENV" == "kind" ]; then
    bash "$SCRIPT_DIR/kind-load-image.sh"
  fi
fi

# ginkgo run --label-filter=basic -v test/e2e
ginkgo -v test/e2e
TEST_EXIT_CODE=$?

echo "Test exit code: $TEST_EXIT_CODE"

if [ "$CREATE_KIND_CLUSTER" == "true" ]; then
  if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "Tests passed, deleting kind cluster: $TEST_CLUSTER_NAME"
    kind delete cluster --name "$TEST_CLUSTER_NAME"
  else
    echo "Tests failed, keeping kind cluster for debugging."
    echo "Kubeconfig: $SAVED_KUBECONFIG"
  fi
fi

exit $TEST_EXIT_CODE
