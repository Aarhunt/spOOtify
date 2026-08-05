import { useState } from 'react';
import { Plus, Check } from 'lucide-react';

interface GenrePillProps {
  name: string;
  onSelect?: (genre: string, isSelected: boolean) => void;
  defaultSelected?: boolean;
}

export function GenrePill({ name, onSelect, defaultSelected = false }: GenrePillProps) {
  const [isSelected, setIsSelected] = useState(defaultSelected);

  const handleClick = () => {
    const nextState = !isSelected;
    setIsSelected(nextState);
    if (onSelect) {
      onSelect(name, nextState);
    }
  };

  return (
    <button
      type="button"
      onClick={handleClick}
      className={`
        inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-medium
        border transition-all duration-200 ease-in-out cursor-pointer select-none
        ${
          isSelected
            ? 'bg-[#1DB954] text-black border-[#1DB954] hover:bg-[#1ed760] hover:scale-105'
            : 'bg-[#181818] text-gray-300 border-[#282828] hover:border-gray-500 hover:text-white'
        }
      `}
    >
      <span>{name}</span>
      {isSelected ? (
        <Check size={14} className="stroke-[3]" />
      ) : (
        <Plus size={14} className="stroke-[2.5]" />
      )}
    </button>
  );
}
