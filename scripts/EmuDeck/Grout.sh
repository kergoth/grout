#!/usr/bin/env bash
set -euo pipefail

scriptdir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)
exec > >(tee -a "${scriptdir}/Grout/grout.log") 2>&1
CFW=EMUDECK LD_LIBRARY_PATH="${scriptdir}/Grout/lib${LD_LIBRARY_PATH:+:${LD_LIBRARY_PATH}}" exec "${scriptdir}/Grout/grout" "$@"
