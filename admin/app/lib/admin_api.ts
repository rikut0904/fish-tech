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
  const response = await fetch(`${API_BASE_URL}/fishes`, { cache: "no-store" });
  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }

  const data = (await response.json()) as ListResponse<Fish>;
  return data.items;
}

export async function fetchPairs(): Promise<FishPair[]> {
  const response = await fetch(`${API_BASE_URL}/pairs`, { cache: "no-store" });
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
  const response = await fetch(`${API_BASE_URL}/admin/fishes`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }

  return (await response.json()) as Fish;
}

export async function uploadFishImage(file: File): Promise<string> {
  const formData = new FormData();
  formData.append("file", file);

  const response = await fetch(`${API_BASE_URL}/admin/fishes/upload-image`, {
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
  const response = await fetch(`${API_BASE_URL}/admin/fishes/${id}`, {
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
  const response = await fetch(`${API_BASE_URL}/admin/pairs`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }

  return (await response.json()) as FishPair;
}

export async function deletePair(id: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/admin/pairs/${id}`, {
    method: "DELETE",
  });
  if (!response.ok) {
    throw new Error(await parseErrorMessage(response));
  }
}
