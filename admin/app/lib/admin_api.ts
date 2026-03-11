export type Fish = {
  id: string;
  name: string;
  category: string;
  description: string;
  imageUrl: string;
  linkUrl: string;
};

export type FishPair = {
  id: string;
  fishIdA: string;
  fishIdB: string;
  score: number;
  memo: string;
};

type ListResponse<T> = {
  items: T[];
};

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";

async function fetchApi(input: string, init?: RequestInit): Promise<Response> {
  try {
    return await fetch(input, init);
  } catch (error) {
    if (error instanceof TypeError) {
      throw new Error(
        "バックエンドへ接続できませんでした。backend が起動しているか確認してください（例: docker compose up -d backend）。",
      );
    }
    throw error;
  }
}

export function buildPairKey(fishIdA: string, fishIdB: string): string {
  return [fishIdA, fishIdB].sort().join(":");
}

export async function parseErrorMessage(response: Response): Promise<string> {
  try {
    const data = (await response.json()) as { error?: string };
    return data.error ?? "APIエラーが発生しました";
  } catch {
    return "APIエラーが発生しました";
  }
}

export async function fetchFishes(): Promise<Fish[]> {
  const response = await fetchApi(`${API_BASE_URL}/fishes`, { cache: "no-store" });
  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }

  const data = (await response.json()) as ListResponse<Fish>;
  return data.items;
}

export async function fetchPairs(): Promise<FishPair[]> {
  const response = await fetchApi(`${API_BASE_URL}/pairs`, { cache: "no-store" });
  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }

  const data = (await response.json()) as ListResponse<FishPair>;
  return data.items;
}

export async function createFish(input: {
  name: string;
  category: string;
  description: string;
  imageUrl: string;
  linkUrl: string;
}): Promise<Fish> {
  const params = new URLSearchParams({
    name: input.name,
    category: input.category,
    description: input.description,
    imageUrl: input.imageUrl,
    linkUrl: input.linkUrl,
  });

  const response = await fetchApi(`${API_BASE_URL}/admin/fishes?${params.toString()}`, {
    method: "POST",
  });
  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }

  return (await response.json()) as Fish;
}

export async function uploadFishImage(file: File): Promise<string> {
  const formData = new FormData();
  formData.append("file", file);

  const response = await fetchApi(`${API_BASE_URL}/admin/fishes/upload-image`, {
    method: "POST",
    body: formData,
  });
  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }

  const data = (await response.json()) as { imageUrl: string };
  return data.imageUrl;
}

export async function deleteFish(id: string): Promise<void> {
  const response = await fetchApi(`${API_BASE_URL}/admin/fishes/${id}`, {
    method: "DELETE",
  });
  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }
}

export async function createPair(input: {
  fishIdA: string;
  fishIdB: string;
  score: number;
  memo: string;
}): Promise<FishPair> {
  const params = new URLSearchParams({
    fishIdA: input.fishIdA,
    fishIdB: input.fishIdB,
    score: String(input.score),
    memo: input.memo,
  });

  const response = await fetchApi(`${API_BASE_URL}/admin/pairs?${params.toString()}`, {
    method: "POST",
  });
  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }

  return (await response.json()) as FishPair;
}

export async function deletePair(id: string): Promise<void> {
  const response = await fetchApi(`${API_BASE_URL}/admin/pairs/${id}`, {
    method: "DELETE",
  });
  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }
}
