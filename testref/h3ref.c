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
        
        LatLng geo;
        setGeoDegs(&geo, lat, lng);
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

    } else if (strcmp(command, "coordijk_distance") == 0) {
        if (argc < 8) {
            fprintf(stderr, "Usage: %s coordijk_distance <i1> <j1> <k1> <i2> <j2> <k2>\n", argv[0]);
            return 1;
        }
        CoordIJK a = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        CoordIJK b = { atoi(argv[5]), atoi(argv[6]), atoi(argv[7]) };
        int d = ijkDistance(&a, &b);
        printf("%d\n", d);

    } else if (strcmp(command, "coordijk_rotate") == 0) {
        if (argc < 6) {
            fprintf(stderr, "Usage: %s coordijk_rotate <ccw|cw> <i> <j> <k>\n", argv[0]);
            return 1;
        }
        const char* which = argv[2];
        CoordIJK v = { atoi(argv[3]), atoi(argv[4]), atoi(argv[5]) };
        if (strcmp(which, "ccw") == 0) {
            _ijkRotate60ccw(&v);
        } else if (strcmp(which, "cw") == 0) {
            _ijkRotate60cw(&v);
        } else {
            fprintf(stderr, "Unknown rotation: %s (use ccw or cw)\n", which);
            return 1;
        }
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "coordijk_hex2d") == 0) {
        if (argc < 5) {
            fprintf(stderr, "Usage: %s coordijk_hex2d <i> <j> <k>\n", argv[0]);
            return 1;
        }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        Vec2d out;
        _ijkToHex2d(&v, &out);
        printf("%.15f %.15f\n", out.x, out.y);

    } else if (strcmp(command, "coordijk_from_hex2d") == 0) {
        if (argc < 4) {
            fprintf(stderr, "Usage: %s coordijk_from_hex2d <x> <y>\n", argv[0]);
            return 1;
        }
        Vec2d in = { atof(argv[2]), atof(argv[3]) };
        CoordIJK out;
        _hex2dToCoordIJK(&in, &out);
        printf("%d %d %d\n", out.i, out.j, out.k);

    } else if (strcmp(command, "coordijk_neighbor") == 0) {
        if (argc < 6) {
            fprintf(stderr, "Usage: %s coordijk_neighbor <digit> <i> <j> <k>\n", argv[0]);
            return 1;
        }
        Direction d = (Direction)atoi(argv[2]);
        CoordIJK v = { atoi(argv[3]), atoi(argv[4]), atoi(argv[5]) };
        _neighbor(&v, d);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "coordijk_up_ap7") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s coordijk_up_ap7 <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _upAp7(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "coordijk_up_ap7r") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s coordijk_up_ap7r <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _upAp7r(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "coordijk_down_ap7") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s coordijk_down_ap7 <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _downAp7(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "coordijk_down_ap7r") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s coordijk_down_ap7r <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _downAp7r(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "coordijk_down_ap3") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s coordijk_down_ap3 <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _downAp3(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "coordijk_down_ap3r") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s coordijk_down_ap3r <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _downAp3r(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "geotofaceijk") == 0) {
        if (argc < 5) {
            fprintf(stderr, "Usage: %s geotofaceijk <lat> <lng> <res>\n", argv[0]);
            return 1;
        }
        double lat = atof(argv[2]);
        double lng = atof(argv[3]);
        int res = atoi(argv[4]);
        LatLng g; setGeoDegs(&g, lat, lng);
        FaceIJK fijk;
        _geoToFaceIjk(&g, res, &fijk);
        printf("%d %d %d %d\n", fijk.face, fijk.coord.i, fijk.coord.j, fijk.coord.k);

    } else if (strcmp(command, "geohex2d") == 0) {
        if (argc < 5) {
            fprintf(stderr, "Usage: %s geohex2d <lat> <lng> <res>\n", argv[0]);
            return 1;
        }
        double lat = atof(argv[2]);
        double lng = atof(argv[3]);
        int res = atoi(argv[4]);
        LatLng g; setGeoDegs(&g, lat, lng);
        int face;
        Vec2d v;
        _geoToHex2d(&g, res, &face, &v);
        printf("%d %.17g %.17g\n", face, v.x, v.y);

    } else {
        fprintf(stderr, "Unknown command: %s\n", command);
        return 1;
    }
    
    return 0;
}
