export type Fish = {
  id: number;
  name: string;
  scientificName: string;
  image: string;
  isFavorite: boolean;
  details: {
    alias: string;
    habitat: string;
    length: string;
    season: string; // 旬
    type: string;   // 種類
    fishingMethod: string;
    recipe: string;
    compatibilityGood: string[];
    compatibilityBad: string[];
  };
};