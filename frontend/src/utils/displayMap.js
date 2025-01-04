import { MapGeoJSON } from "@/data/MapGeo";
import pointInPolygon from "point-in-polygon";

export function IsPointInAustralia(lon, lat) {
    for (const feature of MapGeoJSON.features) {
        if (feature.geometry.type === "Polygon") {
            if (pointInPolygon([lon, lat], feature.geometry.coordinates[0])) {
                return true;
            }
        } else if (feature.geometry.type === "MultiPolygon") {
            for (const polygon of feature.geometry.coordinates) {
                if (pointInPolygon([lon, lat], polygon[0])) {
                    return true;
                }
            }
        }
    }
    return false;
}

export function GetRandomCoordInAustralia() {
    const minLon = 113.15957061;
    const maxLon = 153.61194445;
    const minLat = -43.6345972634;
    const maxLat = -10.6681857235;

    let lon, lat;
    do {
        lon = minLon + Math.random() * (maxLon - minLon);
        lat = minLat + Math.random() * (maxLat - minLat);
    } while (!IsPointInAustralia(lon, lat));

    return { latitude: lat, longitude: lon };
}
