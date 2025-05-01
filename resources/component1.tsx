import React, { useState, useEffect } from 'react';

interface ButtonProps {
  text: string;
  onClick?: () => void;
  variant?: 'primary' | 'secondary' | 'danger';
  disabled?: boolean;
  className?: string;
}

/**
 * A reusable button component with different variants
 */
const Button: React.FC<ButtonProps> = ({
                                         text,
                                         onClick,
                                         variant = 'primary',
                                         disabled = false,
                                         className = '',
                                       }) => {
  const [isHovered, setIsHovered] = useState(false);

  const baseStyles = 'px-4 py-2 rounded font-medium transition-colors';

  const variantStyles = {
    primary: 'bg-blue-500 hover:bg-blue-600 text-white',
    secondary: 'bg-gray-200 hover:bg-gray-300 text-gray-800',
    danger: 'bg-red-500 hover:bg-red-600 text-white',
  };

  const buttonStyles = `${baseStyles} ${variantStyles[variant]} ${className} ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`;

  useEffect(() => {
    // Log when component mounts
    console.log('Button component mounted');
    return () => {
      // Clean up when component unmounts
      console.log('Button component will unmount');
    };
  }, []);

  const handleClick = () => {
    if (!disabled && onClick) {
      onClick();
    }
  };

  return (
      <button
          className={buttonStyles}
          onClick={handleClick}
          disabled={disabled}
          onMouseEnter={() => setIsHovered(true)}
          onMouseLeave={() => setIsHovered(false)}
          data-testid="custom-button"
      >
        {text}
        {isHovered && <span className="ml-1">→</span>}
      </button>
  );
};

export default Button;