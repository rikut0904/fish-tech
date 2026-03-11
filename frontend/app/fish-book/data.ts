export type Fish = {
    id: number;
    name: string;
    scientificName: string;
    details: {
        alias: string;
        habitat: string;
        length: string;
        season: string;
        type: string;
        fishingMethod: string;
        goodCompatibility: string[];
        badCompatibility: string[];
    };
};

export const allFishData: Fish[] = [
    {
        id: 1,
        name: "サンプル魚",
        scientificName: "Sampleus piscis",
        details: {
            alias: "テスト魚",
            habitat: "沿岸",
            length: "30cm",
            season: "通年",
            type: "白身",
            fishingMethod: "定置網",
            goodCompatibility: ["魚A", "魚B"],
            badCompatibility: ["魚C"]
        }
    }
];
