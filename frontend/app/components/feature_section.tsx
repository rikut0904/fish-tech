// サービス特徴紹介
export default function FeatureSection() {
    return (
        <section id="features" className="py-12 bg-white">
            <div className="container mx-auto px-4 grid grid-cols-1 md:grid-cols-3 gap-8">
                <div className="bg-blue-100 rounded-lg p-6 text-center">
                    <h3 className="text-lg font-bold mb-2 text-blue-700">多様な魚種</h3>
                    <p>金沢であまり知られていない魚も多数掲載</p>
                </div>
                <div className="bg-blue-100 rounded-lg p-6 text-center">
                    <h3 className="text-lg font-bold mb-2 text-blue-700">旬の情報</h3>
                    <p>季節ごとのおすすめや食べ方を紹介</p>
                </div>
                <div className="bg-blue-100 rounded-lg p-6 text-center">
                    <h3 className="text-lg font-bold mb-2 text-blue-700">地元の声</h3>
                    <p>漁師や料理人のコメントも掲載予定</p>
                </div>
            </div>
        </section>
    );
}
