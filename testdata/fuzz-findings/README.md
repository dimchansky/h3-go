# Preserved fuzz findings (not auto-executed)

Go fuzz corpus entries that document confirmed findings but are kept **out**
of `testdata/fuzz/` on purpose: entries there are replayed by every ordinary
`go test` run, and the inputs preserved here are pathologically slow, which
would tax every test run for no assertion value. Each entry stays until its
finding is resolved, then moves into a targeted regression test.

## FuzzUpstreamPolygonOperations/bcefd1d03714b6b6

A finite but pathologically slow (~15 s on an Apple M1 Max, go1.26.5)
polygon input discovered by a 30-second fuzz run of
`FuzzUpstreamPolygonOperations` — tracked in
[issue #3](https://github.com/dimchansky/h3-go/issues/3). Because Go's fuzz
engine kills workers that exceed its per-execution budget, this keeps the
target out of the Nightly `-fuzz` rotation (seed-corpus only).

**Root cause (established 2026-07-17, issue #3):** the input decodes to a
16-vertex loop whose latitudes reach ~1e287 radians. All of the time is
burned in `maxPolygonToCellsSizeExperimental`: its rough bounding-box area
estimate divides by `cos(min(|north|, |south|))`, which is *negative* for
these huge angles, so the negative "area" defeats the resolution-coarsening
loop and the size estimate scans every res-6 cell on the planet (~14M
cells) against the 16-vertex loop. The classic `maxPolygonToCellsSize`
takes microseconds on the same input; `polygonToCells*` never run (the
harness caps `size <= 10000`).

**Parity with H3 C is established:** upstream H3 C 4.4.0 exhibits the same
pathology on this input and is slower (~27 s vs ~15 s Go on the same
machine), returning the identical size (7068476). Upstream's own
`fuzzerPolygonToCellsExperimental` uses the same unguarded input domain and
calls the estimator unconditionally, so this is an upstream-reportable
timeout finding (OSS-Fuzz's default single-input budget is 25 s), not a
port defect — no local guard or bound is justified by the parity contract.
The upstream report is tracked in
[docs/FUTURE_WORK.md](../../docs/FUTURE_WORK.md). Cheap C-verified
regression coverage for this input class (NaN/Inf/huge-magnitude
coordinates, including the negative-rough-area mechanism at low
resolutions) lives in `polyfill_maxPolygonToCellsSizeExperimental_test.go`
and `polygon__lineCrossesLine_test.go`; this reproducer stays preserved
here for manual replay because a ~15 s test earns no extra assertion
value.

SHA-256 (also the basis of the Go corpus entry name):

```
bcefd1d03714b6b6658c1649fbe49951fa2ebf91a4cc132c067aa7253168b967
```

Replay (from the repository root):

```sh
mkdir -p testdata/fuzz/FuzzUpstreamPolygonOperations
cp testdata/fuzz-findings/FuzzUpstreamPolygonOperations/bcefd1d03714b6b6 \
   testdata/fuzz/FuzzUpstreamPolygonOperations/
CGO_ENABLED=0 go test -run 'FuzzUpstreamPolygonOperations/bcefd1d03714b6b6' -v .
# --- PASS (finite), but expect roughly 15 s of CPU for this single input.
# Clean up afterwards so ordinary test runs stay fast — remove ONLY the
# copied reproducer (never `rm -r testdata/fuzz`, which would also delete
# any fuzz corpus entries of your own):
rm testdata/fuzz/FuzzUpstreamPolygonOperations/bcefd1d03714b6b6
rmdir testdata/fuzz/FuzzUpstreamPolygonOperations testdata/fuzz 2>/dev/null || true
```
