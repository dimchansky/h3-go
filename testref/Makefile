# H3 Reference Oracle Makefile
# Downloads and builds against H3 C library for test validation

# To upgrade H3 version, simply change this line and run 'make clean-all && make'
H3_VERSION = 4.5.0
H3_URL = https://github.com/uber/h3/archive/v$(H3_VERSION).tar.gz
H3_ARCHIVE = h3-$(H3_VERSION).tar.gz
H3_DIR = h3-$(H3_VERSION)
H3_SRC = $(H3_DIR)/src/h3lib
H3_LIB = $(H3_SRC)/lib
H3_INC = $(H3_SRC)/include

CC = gcc
CFLAGS = -std=c99 -Wall -O2 -I$(H3_INC)
LDFLAGS = -lm

# h3ref command gating (docs/sync/4.4.0-to-4.5.0.md §15.1): the
# geoToFaceIjk/geoToHex2d commands call internal functions that were
# deleted in H3 4.5.0 (_geoToFaceIjk/_geoToHex2d), so they exist only in
# a 4.4.x-configured oracle. At 4.5.0 they fall through to the "Unknown
# command" error.
ifeq ($(H3_VERSION),4.4.0)
CFLAGS += -DH3REF_HAVE_GEO_INTERMEDIATES
endif

.PHONY: all clean download h3-source test version

all: h3ref

version:
	@echo "H3 version: $(H3_VERSION)"
	@echo "To upgrade: Edit H3_VERSION in Makefile, then run 'make clean-all && make'"

# Download and extract H3 source if not present
h3-source: $(H3_DIR)

$(H3_DIR): $(H3_ARCHIVE)
	@echo "Extracting H3 v$(H3_VERSION) source..."
	tar -xzf $(H3_ARCHIVE)
	# Create h3api.h from template
	cp $(H3_INC)/h3api.h.in $(H3_INC)/h3api.h
	@echo "H3 source ready at $(H3_DIR)/"

$(H3_ARCHIVE):
	@echo "Downloading H3 v$(H3_VERSION) from $(H3_URL)..."
	curl -L -o $(H3_ARCHIVE) $(H3_URL)

download: $(H3_ARCHIVE)

# The upstream sources are compiled at recipe time with a shell-expanded
# glob, AFTER the h3-source prerequisite has downloaded and extracted the
# tree: a make-time $(wildcard) would expand while the prerequisite graph
# is built — before extraction on a fresh clone — and yield an empty list.
# The glob is version-derived by construction: it drops the files deleted
# in 4.5.0 (coordijk.c, vec3d.c, vertexGraph.c) and picks up the new ones
# (area.c, cellsToMultiPoly.c). One-shot compilation replaces the old
# per-file object graph; incremental .o reuse is not worth keeping for a
# one-shot oracle build. (docs/sync/4.4.0-to-4.5.0.md §15.1)
h3ref: h3-source h3ref.c
	$(CC) $(CFLAGS) h3ref.c $(H3_LIB)/*.c -o $@ $(LDFLAGS)

clean:
	rm -f h3ref h3ref.o *.o $(H3_LIB)/*.o

clean-all: clean
	rm -rf $(H3_DIR) $(H3_ARCHIVE)

test: h3ref
	@echo "Testing H3 oracle CLI..."
	./h3ref isBaseCellPentagon 4  # Should print 1 (pentagon)
	./h3ref isBaseCellPentagon 0  # Should print 0 (hexagon)
	./h3ref faceIjkToH3 0 0 0 0 0  # Should generate a base H3 index
	@echo "Oracle tests passed!"
