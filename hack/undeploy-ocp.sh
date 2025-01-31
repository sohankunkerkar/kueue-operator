#!/bin/bash

oc delete -f deploy/examples/job.yaml
oc delete -f deploy/crd/
oc delete -f deploy/