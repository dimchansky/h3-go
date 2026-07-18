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
        fprintf(stderr, "  faceIjkToH3 <face> <i> <j> <k> <res> - _faceIjkToH3\n");
        fprintf(stderr, "  latLngToCell <lat> <lng> <res> - latLngToCell\n");
        fprintf(stderr, "  isBaseCellPentagon <baseCell> - _isBaseCellPentagon\n");
        fprintf(stderr, "  h3Rotate60cw <h3index> - _h3Rotate60cw\n");
        fprintf(stderr, "  h3Rotate60ccw <h3index> - _h3Rotate60ccw\n");
        fprintf(stderr, "  getResolution <h3index> - getResolution\n");
        fprintf(stderr, "  getBaseCellNumber <h3index> - getBaseCellNumber\n");
        fprintf(stderr, "  ijkDistance <i1> <j1> <k1> <i2> <j2> <k2> - ijkDistance\n");
        fprintf(stderr, "  ijkRotate60ccw <i> <j> <k> - _ijkRotate60ccw\n");
        fprintf(stderr, "  ijkRotate60cw <i> <j> <k> - _ijkRotate60cw\n");
        fprintf(stderr, "  ijkToHex2d <i> <j> <k> - _ijkToHex2d\n");
        fprintf(stderr, "  hex2dToCoordIJK <x> <y> - _hex2dToCoordIJK\n");
        fprintf(stderr, "  neighbor <digit> <i> <j> <k> - _neighbor\n");
        fprintf(stderr, "  upAp7 <i> <j> <k> - _upAp7\n");
        fprintf(stderr, "  upAp7r <i> <j> <k> - _upAp7r\n");
        fprintf(stderr, "  downAp7 <i> <j> <k> - _downAp7\n");
        fprintf(stderr, "  downAp7r <i> <j> <k> - _downAp7r\n");
        fprintf(stderr, "  downAp3 <i> <j> <k> - _downAp3\n");
        fprintf(stderr, "  downAp3r <i> <j> <k> - _downAp3r\n");
#ifdef H3REF_HAVE_GEO_INTERMEDIATES
        fprintf(stderr, "  geoToFaceIjk <lat> <lng> <res> - _geoToFaceIjk\n");
        fprintf(stderr, "  geoToHex2d <lat> <lng> <res> - _geoToHex2d\n");
#endif
        return 1;
    }

    const char* command = argv[1];
    
    if (strcmp(command, "faceIjkToH3") == 0) {
        if (argc < 7) {
            fprintf(stderr, "Usage: %s faceIjkToH3 <face> <i> <j> <k> <res>\n", argv[0]);
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
        
    } else if (strcmp(command, "latLngToCell") == 0) {
        if (argc < 5) {
            fprintf(stderr, "Usage: %s latLngToCell <lat> <lng> <res>\n", argv[0]);
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
        
    } else if (strcmp(command, "isBaseCellPentagon") == 0) {
        if (argc < 3) {
            fprintf(stderr, "Usage: %s isBaseCellPentagon <baseCell>\n", argv[0]);
            return 1;
        }
        
        int baseCell = atoi(argv[2]);
        int isPent = _isBaseCellPentagon(baseCell);
        printf("%d\n", isPent);
        
    } else if (strcmp(command, "h3Rotate60cw") == 0) {
        if (argc < 3) {
            fprintf(stderr, "Usage: %s h3Rotate60cw <h3index>\n", argv[0]);
            return 1;
        }
        
        H3Index h3 = (H3Index)strtoull(argv[2], NULL, 16);
        H3Index rotated = _h3Rotate60cw(h3);
        printf("0x%" PRIx64 "\n", rotated);
        
    } else if (strcmp(command, "h3Rotate60ccw") == 0) {
        if (argc < 3) {
            fprintf(stderr, "Usage: %s h3Rotate60ccw <h3index>\n", argv[0]);
            return 1;
        }
        
        H3Index h3 = (H3Index)strtoull(argv[2], NULL, 16);
        H3Index rotated = _h3Rotate60ccw(h3);
        printf("0x%" PRIx64 "\n", rotated);
        
    } else if (strcmp(command, "getResolution") == 0) {
        if (argc < 3) {
            fprintf(stderr, "Usage: %s getResolution <h3index>\n", argv[0]);
            return 1;
        }
        
        H3Index h3 = (H3Index)strtoull(argv[2], NULL, 16);
        int res = getResolution(h3);
        printf("%d\n", res);
        
    } else if (strcmp(command, "getBaseCellNumber") == 0) {
        if (argc < 3) {
            fprintf(stderr, "Usage: %s getBaseCellNumber <h3index>\n", argv[0]);
            return 1;
        }
        
        H3Index h3 = (H3Index)strtoull(argv[2], NULL, 16);
        int baseCell = getBaseCellNumber(h3);
        printf("%d\n", baseCell);

    } else if (strcmp(command, "ijkDistance") == 0) {
        if (argc < 8) {
            fprintf(stderr, "Usage: %s ijkDistance <i1> <j1> <k1> <i2> <j2> <k2>\n", argv[0]);
            return 1;
        }
        CoordIJK a = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        CoordIJK b = { atoi(argv[5]), atoi(argv[6]), atoi(argv[7]) };
        int d = ijkDistance(&a, &b);
        printf("%d\n", d);
    
    } else if (strcmp(command, "ijkRotate60ccw") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s ijkRotate60ccw <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _ijkRotate60ccw(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);
    
    } else if (strcmp(command, "ijkRotate60cw") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s ijkRotate60cw <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _ijkRotate60cw(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "ijkToHex2d") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s ijkToHex2d <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        Vec2d out;
        _ijkToHex2d(&v, &out);
        printf("%.15f %.15f\n", out.x, out.y);

    } else if (strcmp(command, "hex2dToCoordIJK") == 0) {
        if (argc < 4) { fprintf(stderr, "Usage: %s hex2dToCoordIJK <x> <y>\n", argv[0]); return 1; }
        Vec2d in = { atof(argv[2]), atof(argv[3]) };
        CoordIJK out;
        _hex2dToCoordIJK(&in, &out);
        printf("%d %d %d\n", out.i, out.j, out.k);

    } else if (strcmp(command, "neighbor") == 0) {
        if (argc < 6) {
            fprintf(stderr, "Usage: %s neighbor <digit> <i> <j> <k>\n", argv[0]);
            return 1;
        }
        Direction d = (Direction)atoi(argv[2]);
        CoordIJK v = { atoi(argv[3]), atoi(argv[4]), atoi(argv[5]) };
        _neighbor(&v, d);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "upAp7") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s upAp7 <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _upAp7(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "upAp7r") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s upAp7r <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _upAp7r(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "downAp7") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s downAp7 <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _downAp7(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "downAp7r") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s downAp7r <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _downAp7r(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "downAp3") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s downAp3 <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _downAp3(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

    } else if (strcmp(command, "downAp3r") == 0) {
        if (argc < 5) { fprintf(stderr, "Usage: %s downAp3r <i> <j> <k>\n", argv[0]); return 1; }
        CoordIJK v = { atoi(argv[2]), atoi(argv[3]), atoi(argv[4]) };
        _downAp3r(&v);
        printf("%d %d %d\n", v.i, v.j, v.k);

#ifdef H3REF_HAVE_GEO_INTERMEDIATES
    /* _geoToFaceIjk/_geoToHex2d exist only in the 4.4.x tree; the 4.5.0
       Vec3 pipeline removed them (docs/sync/4.4.0-to-4.5.0.md §15.1). */
    } else if (strcmp(command, "geoToFaceIjk") == 0) {
        if (argc < 5) {
            fprintf(stderr, "Usage: %s geoToFaceIjk <lat> <lng> <res>\n", argv[0]);
            return 1;
        }
        double lat = atof(argv[2]);
        double lng = atof(argv[3]);
        int res = atoi(argv[4]);
        LatLng g; setGeoDegs(&g, lat, lng);
        FaceIJK fijk;
        _geoToFaceIjk(&g, res, &fijk);
        printf("%d %d %d %d\n", fijk.face, fijk.coord.i, fijk.coord.j, fijk.coord.k);

    } else if (strcmp(command, "geoToHex2d") == 0) {
        if (argc < 5) {
            fprintf(stderr, "Usage: %s geoToHex2d <lat> <lng> <res>\n", argv[0]);
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
#endif /* H3REF_HAVE_GEO_INTERMEDIATES */

    } else {
        fprintf(stderr, "Unknown command: %s\n", command);
        return 1;
    }
    
    return 0;
}
