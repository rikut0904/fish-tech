import React, { Suspense } from 'react';
import RecipeSampleClient from './RecipeSampleClient';

export default function RecipeSamplePage() {
    return (
        <Suspense fallback={<div className="p-6">読み込み中…</div>}>
            {/* クライアント側のインタラクティブ部分を Suspense でラップ */}
            <RecipeSampleClient />
        </Suspense>
    );
}
