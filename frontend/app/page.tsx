import Image from "next/image";
import Header from "./components/header";
import HeroSection from "./components/hero_section";
import FeatureSection from "./components/feature_section";
import FishCardList from "./components/fish_card_list";
import Footer from "./components/footer";

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen bg-blue-50">
      <Header />
      <main className="flex-1">
        <HeroSection />
        <FeatureSection />
        <section className="container mx-auto px-4 py-8">
          <h2 className="text-xl font-bold mb-4">おすすめの魚種</h2>
          <FishCardList />
        </section>
      </main>
      <Footer />
    </div>
  );
}
