#!/usr/bin/env bash
set -euo pipefail

scriptdir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)
exec > >(tee -a "${scriptdir}/grout.log") 2>&1
CFW=RETRODECK LD_LIBRARY_PATH="${scriptdir}/lib${LD_LIBRARY_PATH:+:${LD_LIBRARY_PATH}}" exec "${scriptdir}/grout" "$@"
