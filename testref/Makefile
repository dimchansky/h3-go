# H3 Reference Oracle Makefile
# Downloads and builds against H3 C library for test validation

# To upgrade H3 version, simply change this line and run 'make clean-all && make'
H3_VERSION = 4.3.0
H3_URL = https://github.com/uber/h3/archive/v$(H3_VERSION).tar.gz
H3_ARCHIVE = h3-$(H3_VERSION).tar.gz
H3_DIR = h3-$(H3_VERSION)
H3_SRC = $(H3_DIR)/src/h3lib
H3_LIB = $(H3_SRC)/lib
H3_INC = $(H3_SRC)/include

CC = gcc
CFLAGS = -std=c99 -Wall -O2 -I$(H3_INC)
LDFLAGS = -lm

# H3 source files needed for the oracle
H3_SOURCES = \
    $(H3_LIB)/algos.c \
    $(H3_LIB)/baseCells.c \
    $(H3_LIB)/bbox.c \
    $(H3_LIB)/coordijk.c \
    $(H3_LIB)/directedEdge.c \
    $(H3_LIB)/faceijk.c \
    $(H3_LIB)/h3Assert.c \
    $(H3_LIB)/h3Index.c \
    $(H3_LIB)/iterators.c \
    $(H3_LIB)/latLng.c \
    $(H3_LIB)/linkedGeo.c \
    $(H3_LIB)/localij.c \
    $(H3_LIB)/mathExtensions.c \
    $(H3_LIB)/polyfill.c \
    $(H3_LIB)/polygon.c \
    $(H3_LIB)/vec2d.c \
    $(H3_LIB)/vec3d.c \
    $(H3_LIB)/vertex.c \
    $(H3_LIB)/vertexGraph.c

OBJECTS = $(H3_SOURCES:.c=.o)

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

h3ref: h3-source h3ref.o $(OBJECTS)
	$(CC) -o $@ h3ref.o $(OBJECTS) $(LDFLAGS)

h3ref.o: h3ref.c $(H3_DIR)
	$(CC) $(CFLAGS) -c -o $@ $<

%.o: %.c $(H3_DIR)
	$(CC) $(CFLAGS) -c -o $@ $<

clean:
	rm -f h3ref h3ref.o $(OBJECTS)

clean-all: clean
	rm -rf $(H3_DIR) $(H3_ARCHIVE)

test: h3ref
	@echo "Testing H3 oracle CLI..."
	./h3ref pentagon 4  # Should print 1 (pentagon)
	./h3ref pentagon 0  # Should print 0 (hexagon) 
	./h3ref faceijk 0 0 0 0 0  # Should generate a base H3 index
	@echo "Oracle tests passed!"