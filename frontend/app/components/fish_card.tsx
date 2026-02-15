// 個別魚カード
interface FishCardProps {
    name: string;
    desc: string;
    image: string;
}

export default function FishCard({ name, desc, image }: FishCardProps) {
    return (
        <div className="bg-white rounded-lg shadow-md p-4 flex flex-col items-center">
            <img src={image} alt={name} className="w-32 h-24 object-contain mb-2" />
            <h4 className="text-lg font-bold text-blue-700 mb-1">{name}</h4>
            <p className="text-sm text-gray-700 text-center">{desc}</p>
        </div>
    );
}
