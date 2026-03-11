"use client";

import { useReducer } from "react";

type PatchAction<T> = Partial<T> | ((prev: T) => Partial<T>);

function patchReducer<T>(state: T, action: PatchAction<T>): T {
  const patch = typeof action === "function" ? action(state) : action;
  return { ...state, ...patch };
}

export function usePatchReducer<T>(initialState: T) {
  return useReducer(patchReducer<T>, initialState);
}
