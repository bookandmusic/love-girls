import axios from "axios";

export interface GeocodingResult {
  lat: string;
  lon: string;
  displayName: string;
}

const nominatimApi = axios.create({
  baseURL: "https://nominatim.openstreetmap.org",
  timeout: 15000,
  headers: {
    "Content-Type": "application/json",
  },
});

export const geocodingApi = {
  async search(query: string): Promise<GeocodingResult[]> {
    if (!query.trim()) {
      return [];
    }

    try {
      const response = await nominatimApi.get("/search", {
        params: {
          q: query,
          format: "json",
          limit: 5,
          "accept-language": "zh",
        },
      });

      return response.data.map(
        (item: { lat: string; lon: string; display_name: string }) => ({
          lat: item.lat,
          lon: item.lon,
          displayName: item.display_name,
        }),
      );
    } catch (error) {
      console.error("Geocoding search error:", error);
      throw error;
    }
  },
};
