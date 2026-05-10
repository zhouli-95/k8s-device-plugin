#!/bin/bash

KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"


kind load docker-image device-plugin:dev --name "$KIND_CLUSTER_NAME"

kind load docker-image ubuntu:24.04 --name "$KIND_CLUSTER_NAME"