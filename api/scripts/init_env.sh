#!/usr/bin/env bash
set -e

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

CUR_DIR=$(pwd)

source $CUR_DIR/.env.local