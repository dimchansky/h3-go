#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <inttypes.h>
#include "h3api.h"
#include "latLng.h"
#include "faceijk.h"
#include "coordijk.h"
#include "baseCells.h"
#include "h3Index.h"

int main(int argc, char *argv[]) {
    if (argc < 2) {
        fprintf(stderr, "Usage: %s <command> [args...]\n", argv[0]);
        fprintf(stderr, "Commands:\n");
        fprintf(stderr, "  faceijk <face> <i> <j> <k> <res> - Convert FaceIJK to H3 index\n");
        fprintf(stderr, "  latlng <lat> <lng> <res> - Convert LatLng to H3 index\n");
        fprintf(stderr, "  pentagon <baseCell> - Check if base cell is pentagon\n");
        fprintf(stderr, "  rotate60cw <h3index> - Rotate H3 index 60 degrees clockwise\n");
        fprintf(stderr, "  rotate60ccw <h3index> - Rotate H3 index 60 degrees counter-clockwise\n");
        fprintf(stderr, "  resolution <h3index> - Get H3 index resolution\n");
        fprintf(stderr, "  basecell <h3index> - Get H3 index base cell\n");
        return 1;
    }

    const char* command = argv[1];
    
    if (strcmp(command, "faceijk") == 0) {
        if (argc < 7) {
            fprintf(stderr, "Usage: %s faceijk <face> <i> <j> <k> <res>\n", argv[0]);
            return 1;
        }
        
        int face = atoi(argv[2]);
        int i = atoi(argv[3]);
        int j = atoi(argv[4]);
        int k = atoi(argv[5]);
        int res = atoi(argv[6]);
        
        FaceIJK fijk = {face, {i, j, k}};
        H3Index h3 = _faceIjkToH3(&fijk, res);
        
        printf("0x%" PRIx64 "\n", h3);
        
    } else if (strcmp(command, "latlng") == 0) {
        if (argc < 5) {
            fprintf(stderr, "Usage: %s latlng <lat> <lng> <res>\n", argv[0]);
            return 1;
        }
        
        double lat = atof(argv[2]);
        double lng = atof(argv[3]);
        int res = atoi(argv[4]);
        
        LatLng geo = {lat, lng};
        H3Index h3;
        H3Error err = latLngToCell(&geo, res, &h3);
        
        if (err != E_SUCCESS) {
            printf("0x0 %d\n", (int)err);
        } else {
            printf("0x%" PRIx64 " 0\n", h3);
        }
        
    } else if (strcmp(command, "pentagon") == 0) {
        if (argc < 3) {
            fprintf(stderr, "Usage: %s pentagon <baseCell>\n", argv[0]);
            return 1;
        }
        
        int baseCell = atoi(argv[2]);
        int isPent = _isBaseCellPentagon(baseCell);
        printf("%d\n", isPent);
        
    } else if (strcmp(command, "rotate60cw") == 0) {
        if (argc < 3) {
            fprintf(stderr, "Usage: %s rotate60cw <h3index>\n", argv[0]);
            return 1;
        }
        
        H3Index h3 = (H3Index)strtoull(argv[2], NULL, 16);
        H3Index rotated = _h3Rotate60cw(h3);
        printf("0x%" PRIx64 "\n", rotated);
        
    } else if (strcmp(command, "rotate60ccw") == 0) {
        if (argc < 3) {
            fprintf(stderr, "Usage: %s rotate60ccw <h3index>\n", argv[0]);
            return 1;
        }
        
        H3Index h3 = (H3Index)strtoull(argv[2], NULL, 16);
        H3Index rotated = _h3Rotate60ccw(h3);
        printf("0x%" PRIx64 "\n", rotated);
        
    } else if (strcmp(command, "resolution") == 0) {
        if (argc < 3) {
            fprintf(stderr, "Usage: %s resolution <h3index>\n", argv[0]);
            return 1;
        }
        
        H3Index h3 = (H3Index)strtoull(argv[2], NULL, 16);
        int res = getResolution(h3);
        printf("%d\n", res);
        
    } else if (strcmp(command, "basecell") == 0) {
        if (argc < 3) {
            fprintf(stderr, "Usage: %s basecell <h3index>\n", argv[0]);
            return 1;
        }
        
        H3Index h3 = (H3Index)strtoull(argv[2], NULL, 16);
        int baseCell = getBaseCellNumber(h3);
        printf("%d\n", baseCell);
        
    } else {
        fprintf(stderr, "Unknown command: %s\n", command);
        return 1;
    }
    
    return 0;
}