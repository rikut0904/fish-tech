export interface Fish {
  id: number;
  name: string;
  scientificName: string;
  image: string;
  isFavorite: boolean;
  details: {
    alias: string;         // 別名
    habitat: string;       // 生息地
    length: string;        // 体長
    season: string;        // 旬
    type: string;          // 種類
    fishingMethod: string; // 漁法
    recipe: string;        // おすすめレシピ
    goodCompatibility: string[]; // 相性の良い食材
    badCompatibility: string[];  // 相性の悪い食材
  };
}