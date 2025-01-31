#!/bin/bash

oc delete -f deploy/crd/
oc delete -f deploy/
oc delete -f deploy/examples/job.yaml