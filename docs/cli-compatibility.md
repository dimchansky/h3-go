# `h3` command-line compatibility

This is the authoritative compatibility contract for the pure-Go `h3`
executable; the user-facing quick start is [cmd/h3/README.md](../cmd/h3/README.md).

The authoritative H3 C v4.4.0 build definition creates target `h3_bin` and
sets `OUTPUT_NAME h3`. This repository therefore installs the compatible Go
executable as **`h3`**.

## Install and architecture

```sh
go install github.com/dimchansky/h3-go/cmd/h3@latest
h3 --help
```

`cmd/h3/main.go` only passes process arguments and streams to
`internal/cli.Run`. The internal package owns a small dependency-free parser,
command groups, input scanners, and renderers. `Run` is stateless, accepts
injected stdin/stdout/stderr, and returns an exit code, so it is safe for
repeated in-process tests. Operations use the public `h3` package API; no CLI
exports, cgo, `unsafe`, or external dependencies were required.

## Commands and formats

All 63 upstream commands are implemented:

- indexing and inspection: `cellToLatLng`, `latLngToCell`,
  `cellToBoundary`, `getResolution`, `getBaseCellNumber`, `getIndexDigit`,
  `constructCell`, `stringToInt`, `intToString`, `isValidCell`,
  `isResClassIII`, `isPentagon`, `getIcosahedronFaces`;
- traversal and hierarchy: `gridDisk`, `gridDiskDistances`, `gridRing`,
  `gridPathCells`, `gridDistance`, `cellToLocalIj`, `localIjToCell`,
  `cellToParent`, `cellToChildren`, `cellToChildrenSize`,
  `cellToCenterChild`, `cellToChildPos`, `childPosToCell`, `compactCells`,
  `uncompactCells`;
- regions: `polygonToCells`, `maxPolygonToCellsSize`,
  `cellsToMultiPolygon`;
- directed edges and vertexes: `areNeighborCells`, `cellsToDirectedEdge`,
  `isValidDirectedEdge`, `getDirectedEdgeOrigin`,
  `getDirectedEdgeDestination`, `directedEdgeToCells`,
  `originToDirectedEdges`, `directedEdgeToBoundary`, `cellToVertex`,
  `cellToVertexes`, `vertexToLatLng`, `isValidVertex`;
- measurements and utilities: `degsToRads`, `radsToDegs`,
  `getHexagonAreaAvgKm2`, `getHexagonAreaAvgM2`, `cellAreaRads2`,
  `cellAreaKm2`, `cellAreaM2`, `getHexagonEdgeLengthAvgKm`,
  `getHexagonEdgeLengthAvgM`, `edgeLengthRads`, `edgeLengthKm`,
  `edgeLengthM`, `getNumCells`, `getRes0Cells`, `getPentagons`,
  `pentagonCount`, `greatCircleDistanceRads`, `greatCircleDistanceKm`,
  `greatCircleDistanceM`, `describeH3Error`.

The exact short/long aliases, requiredness, defaults, and permitted formats
are the reviewed rows in [cli-contract.csv](cli-contract.csv). Supported
renderings are the upstream JSON, WKT, newline, and numeric forms. Polygon and
coordinate inputs are JSON arrays in latitude/longitude order. Cell-set and
polygon commands accept inline input, files, and `-i --` for stdin exactly
where upstream does.

Examples:

```sh
h3 latLngToCell -r 9 --lat 37.7759 --lng -122.4180
h3 cellToBoundary -c 8928308280fffff -f wkt
printf '[[0, 0], [0, 1]]' | h3 greatCircleDistanceKm -i --
h3 compactCells -i cells.txt -f newline
```

## Compatibility policy

- Successful commands exit 0. H3 operation failures exit with the numeric H3
  error code (1–19) and write the classification to stderr.
- To preserve the C parser contract, help and recognized-subcommand argument
  errors exit 0; no command or an unknown command exits 1.
- Stable upstream traversal and serialization order is retained. Differential
  JSON coordinate comparisons allow at most `5e-8` degrees (needed only for
  platform `libm` behavior extremely near a pole); scalar metrics use a
  `1e-12` relative tolerance.
- The upstream 1500-byte cell-input scanner is reproduced, including chunk
  boundary behavior visible in its multipolygon fixture.
- Help whitespace and diagnostic wording are not byte-locked. The additive
  `--version` option reports H3 compatibility plus Go module/build metadata.

## Inventories and tests

- [cli-contract.csv](cli-contract.csv): semantic command/flag/format mapping.
- [cli-test-inventory.csv](cli-test-inventory.csv): all 170 upstream commands,
  expected outputs, source locations, and source hashes.
- [cli-fixture-inventory.csv](cli-fixture-inventory.csv): all 15 referenced
  fixture hashes.
- [cli-source-inventory.csv](cli-source-inventory.csv): CMake, parser, and
  command implementation sources that define the contract.

`make test-cli` executes all 170 cases in-process. `make test-cli-process`
builds the real binary and checks pipes, stderr, and exit statuses.
`make test-cli-diff` builds upstream `h3_bin` and differentially executes the
full registry against both binaries. `make check-cli-inventory` fails if a
command, test, fixture, parser source, build definition, or reviewed hash
changes.

During an upstream upgrade, run the general `upstream-diff`, then the CLI
inventory gate against the new tree. Review and update every changed semantic
row before changing its hash; port new scenarios and run the C differential
suite before accepting the new compatibility version.
