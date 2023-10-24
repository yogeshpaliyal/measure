"use client"

import React, { useEffect, useRef, useState } from 'react';

interface CheckboxDropdownProps {
  title: string,
  items: string[];
}

const CheckboxDropdown: React.FC<CheckboxDropdownProps> = ({ title, items }) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement | null>(null);


  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    const handleFocusIn = (event: FocusEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
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

  const toggleDropdown = () => {
    setIsOpen(!isOpen);
  };

  return (
    <div className="z-50 relative inline-block text-left" ref={dropdownRef} >
      <div>
        <button
          type="button"
          onClick={toggleDropdown}
          className="inline-flex justify-center w-full text-black font-display border border-black rounded-md outline-none hover:bg-yellow-200 focus:bg-yellow-200 active:bg-yellow-300"
        >
          <span className="px-6 py-2">{title}</span>
          <span className="border border-black border-t-0 border-r-0 border-b-0 px-4 py-2">⏷</span>
        </button>
      </div>

      {isOpen && (
        <div className="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg ring-1 ring-black ring-opacity-5">
          <div
            role="menu"
            aria-orientation="vertical"
            aria-labelledby="options-menu"
          >
            {items.map((item) => (
              <div key={item} className="px-2 py-2 bg-neutral-950 text-white hover:text-black hover:bg-yellow-200 font-display text-left">
                <input
                  type="checkbox"
                  className="appearance-none border-white rounded-sm text-black font-display bg-neutral-950 focus:ring-offset-yellow-200 focus:ring-0 checked:ring-1 checked:ring-white"
                  value={item}
                />
                <span className="ml-2">{item}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default CheckboxDropdown;