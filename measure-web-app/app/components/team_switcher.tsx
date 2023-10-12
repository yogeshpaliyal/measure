"use client"

import React, { useEffect, useRef, useState } from 'react';

interface TeamSwitcherProps {
  items: string[];
}

const TeamSwitcher: React.FC<TeamSwitcherProps> = ({ items }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedItem, setSelectedItem] = useState<string | null>(null);
  const TeamSwitcherRef = useRef<HTMLDivElement | null>(null);


  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        TeamSwitcherRef.current &&
        !TeamSwitcherRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    const handleFocusIn = (event: FocusEvent) => {
      if (
        TeamSwitcherRef.current &&
        !TeamSwitcherRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    document.addEventListener('focusin', handleFocusIn);

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
      document.removeEventListener('focusin', handleFocusIn);
    };
  }, []);

  const toggleTeamSwitcher = () => {
    setIsOpen(!isOpen);
  };

  const selectItem = (item: string) => {
    setSelectedItem(item);
    setIsOpen(false);
  };

  return (
    <div className="z-50 relative w-32 self-center inline-block text-left" ref={TeamSwitcherRef} >
      <div>
        <button
          type="button"
          onClick={toggleTeamSwitcher}
          className="aspect-square w-full text-black text-xl font-display border border-black rounded-full outline-none hover:bg-yellow-200 focus:bg-yellow-200 active:bg-yellow-300"
        >
          <span className="pl-4">{selectedItem ? selectedItem : items[0]}</span>
          <span className="pl-3 pr-4 text-sm">⏷</span>
        </button>
      </div>

      {isOpen && (
        <div className="origin-top-right absolute left-0 mt-2 w-48 rounded-md shadow-lg ring-1 ring-black ring-opacity-5">
          <div
            role="menu"
            aria-orientation="vertical"
            aria-labelledby="options-menu"
          >
            {items.map((item) => (
              <button
                key={item}
                onClick={() => selectItem(item)}
                className="block w-full px-2 py-2 text-white bg-neutral-950 font-display text-left hover:text-black hover:bg-yellow-200 active:bg-yellow-300 outline-none focus:bg-yellow-200"
                role="menuitem"
              >
                {item}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default TeamSwitcher;