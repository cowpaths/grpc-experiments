# grpc-experiments

This repository contains a repro of a grpc bug as described here: https://github.com/grpc/grpc-go/issues/674

The code is heavily based on this gist, referenced from the prior url: https://gist.github.com/akhenakh/adb03aa05ac1a6e85611422128d2a004

## To run:

1. Change the constants in config.go to match your project.
1. If necessary, run `setup.sh` to create a dummy table.
1. Run `run.sh` to do a single test.  I find that (ballbark) ~1 out of 5 attempts will successfully demonstrate the bug, which is indicated by a runtime of several minutes (as opposed to ~10 seconds in a normal run).
