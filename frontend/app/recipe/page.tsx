'use client';

import { useState } from 'react';
import Image from 'next/image';
import { ChefHat, Search, Clock, DollarSign, Star } from 'lucide-react';

type Fish = {
  id: string;
  name: string;
  image: string;
};

type Recipe = {
  id: string;
  title: string;
  imageUrl: string;
  cookingTime?: string;
  cost?: string;
  score?: number;
  explain?: string;
};

const fishes: Fish[] = [
  { id: '1', name: 'サーモン', image: '/recipe/salmon.png' },
  { id: '2', name: 'マグロ', image: '/recipes-sample/image.png' },
  { id: '3', name: 'タイ', image: '/recipe/tai.png' },
];

const recipesData: Record<string, Recipe[]> = {
  '1': [
    {
      id: 'r1',
      title: 'サーモンのムニエル',
      imageUrl: '/recipes/salmon.jpg',
      cookingTime: '15分',
      cost: '500円前後',
      score: 5,
      explain: 'バター醤油で香ばしく焼き上げます。',
    },
  ],
  '2': [
    {
      id: 'r2',
      title: 'マグロの漬け丼',
      imageUrl: '/recipes/tuna.jpg',
      cookingTime: '10分',
      cost: '600円前後',
      score: 4,
      explain: '特製ダレで簡単に作れます。',
    },
  ],
};

export default function RecipePage() {
  const [selectedFishId, setSelectedFishId] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');

  const selectedFish = fishes.find((f) => f.id === selectedFishId);
  const recipes = selectedFishId ? recipesData[selectedFishId] || [] : [];

  const filteredRecipes = recipes.filter((recipe) => {
    const query = searchQuery.toLowerCase();
    return (
      recipe.title.toLowerCase().includes(query) ||
      recipe.explain?.toLowerCase().includes(query)
    );
  });

  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-50 to-white p-6">
      <header className="mb-8 flex items-center gap-3">
        <ChefHat className="text-blue-600 size-8" />
        <h1 className="text-3xl font-bold">魚のレシピ</h1>
      </header>

      {/* 魚選択 */}
      <div className="mb-8">
        <h2 className="text-xl mb-4 font-semibold">魚を選択</h2>

        <div className="grid grid-cols-3 gap-4">
          {fishes.map((fish) => (
            <button
              key={fish.id}
              onClick={() => setSelectedFishId(fish.id)}
              className={`border rounded-lg overflow-hidden ${
                selectedFishId === fish.id
                  ? 'border-blue-600 shadow-lg'
                  : 'border-gray-300'
              }`}
            >
              <div className="relative aspect-square">
                <Image
                  src={fish.image}
                  alt={fish.name}
                  fill
                  className="object-contain"
                />
              </div>
              <div className="p-2 text-center text-sm font-medium">
                {fish.name}
              </div>
            </button>
          ))}
        </div>
      </div>

      {/* レシピ一覧 */}
      {selectedFish && (
        <div>
          <div className="mb-4">
            <h2 className="text-2xl font-bold mb-2">
              {selectedFish.name}のおすすめレシピ
            </h2>

            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 size-4" />
              <input
                type="text"
                placeholder="レシピ検索..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border rounded-lg"
              />
            </div>
          </div>

          {filteredRecipes.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {filteredRecipes.map((recipe) => (
                <div
                  key={recipe.id}
                  className="bg-white rounded-lg shadow-md overflow-hidden"
                >
                  <div className="relative h-48">
                    <Image
                      src={recipe.imageUrl}
                      alt={recipe.title}
                      fill
                      className="object-contain"
                    />
                  </div>

                  <div className="p-4">
                    <h3 className="text-lg font-semibold mb-2">
                      {recipe.title}
                      {recipe.score !== undefined && (
                        <span className="ml-2 text-sm text-gray-500">
                          ({recipe.score}点)
                        </span>
                      )}
                    </h3>

                    <div className="flex gap-4 text-sm text-gray-600 mb-2">
                      {recipe.cookingTime && (
                        <div className="flex items-center gap-1">
                          <Clock size={16} />
                          {recipe.cookingTime}
                        </div>
                      )}

                      {recipe.cost && (
                        <div className="flex items-center gap-1">
                          <DollarSign size={16} />
                          {recipe.cost}
                        </div>
                      )}
                    </div>

                    {recipe.score && (
  <div className="flex gap-1 mb-2">
    {[1, 2, 3, 4, 5].map((star) => {
      const score = recipe.score ?? 0;
      return (
        <Star
          key={star}
          size={16}
          className={
            star <= score
              ? 'text-yellow-400 fill-yellow-400'
              : 'text-gray-300'
          }
        />
      );
    })}
  </div>
)}

                    {recipe.explain && (
                      <p className="text-sm text-gray-700">
                        {recipe.explain}
                      </p>
                    )}
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500 mt-6">
              レシピが見つかりません
            </p>
          )}
        </div>
      )}
    </div>
  );
}