#!/usr/bin/env bash

set -xeuo pipefail

function list-bound-apps() {
  kubectl get --all-namespaces servicebindings.servicebinding.io \
    -o=custom-columns="NAMESPACE":"metadata.namespace","APP_GUID":".metadata.labels.korifi\.cloudfoundry\.org/app-guid" \
    --no-headers | sort | uniq
}

function main() {
  apps="$(list-bound-apps)"
  if [[ -z "${apps}" ]]; then
    echo "No apps bound to services. Nothing to do."
    return
  fi

  while IFS= read -r line; do
    read -r ns app_guid <<<$line

    while ! diff actual desrired; do
      echo "waiting for status.serviceBindings in cfapp $ns/%app_guid"
      sleep 1
      kubectl get --all-namespaces servicebindings.servicebinding.io \
        -l "korifi.cloudfoundry.org/app-guid=$app_guid" \
        -o=custom-columns="NAME":"metadata.name" \
        --no-headers | sort
      >desired
      kubectl -n $ns get cfapps.korifi.cloudfoundry.org $app_guid -o=jsonpath='{$.status.serviceBindings[*].name}' | tr ' ' '\n' | sort >actual
    done

    kubectl delete --all-namespaces servicebindings.servicebinding.io -l "korifi.cloudfoundry.org/app-guid=$app_guid"
    kubectl delete --all-namespaces appworkloads.korifi.cloudfoundry.org -l "korifi.cloudfoundry.org/app-guid=$app_guid"
  done <<<"$apps"

  # # bindings=("<none>")
  # bindings=("<none>")
  # while [ ${bindings[0]} == "<none>" ]; do
  #   echo "in while"
  #   sleep 1
  #   bindings="fdsfdsfsd"
  # done
}

main
