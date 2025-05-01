import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Button from '../component1';

// This file demonstrates how to test React components using React Testing Library
// It includes examples of common testing patterns and best practices

describe('Button Component (RTL Example)', () => {
  // Test that the component renders with the correct text
  test('renders button with correct text', () => {
    render(<Button text="Click me" />);
    
    // Use getByText to find an element with the given text
    const buttonElement = screen.getByText('Click me');
    expect(buttonElement).toBeInTheDocument();
    
    // Alternative: using test id
    const buttonByTestId = screen.getByTestId('custom-button');
    expect(buttonByTestId).toHaveTextContent('Click me');
  });

  // Test different variants of the button
  test('renders button with correct variant styles', () => {
    const { rerender } = render(<Button text="Primary" variant="primary" />);
    
    // Check primary variant
    const primaryButton = screen.getByTestId('custom-button');
    expect(primaryButton.className).toContain('bg-blue-500');
    
    // Re-render with different variant
    rerender(<Button text="Secondary" variant="secondary" />);
    const secondaryButton = screen.getByTestId('custom-button');
    expect(secondaryButton.className).toContain('bg-gray-200');
    
    // Re-render with danger variant
    rerender(<Button text="Danger" variant="danger" />);
    const dangerButton = screen.getByTestId('custom-button');
    expect(dangerButton.className).toContain('bg-red-500');
  });

  // Test the onClick handler
  test('calls onClick handler when clicked', () => {
    const handleClick = jest.fn();
    render(<Button text="Click me" onClick={handleClick} />);
    
    const buttonElement = screen.getByTestId('custom-button');
    userEvent.click(buttonElement);
    
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  // Test disabled state
  test('does not call onClick when disabled', () => {
    const handleClick = jest.fn();
    render(<Button text="Click me" onClick={handleClick} disabled={true} />);
    
    const buttonElement = screen.getByTestId('custom-button');
    userEvent.click(buttonElement);
    
    expect(handleClick).not.toHaveBeenCalled();
    expect(buttonElement).toBeDisabled();
    expect(buttonElement.className).toContain('opacity-50');
  });

  // Test hover state
  test('shows arrow on hover', async () => {
    render(<Button text="Hover me" />);
    
    const buttonElement = screen.getByTestId('custom-button');
    
    // Initially, arrow should not be visible
    expect(buttonElement).not.toHaveTextContent('→');
    
    // Trigger mouse enter
    fireEvent.mouseEnter(buttonElement);
    
    // Arrow should now be visible
    await waitFor(() => {
      expect(buttonElement).toHaveTextContent('→');
    });
    
    // Trigger mouse leave
    fireEvent.mouseLeave(buttonElement);
    
    // Arrow should be hidden again
    await waitFor(() => {
      expect(buttonElement).not.toHaveTextContent('→');
    });
  });

  // Test custom class names
  test('applies custom className', () => {
    render(<Button text="Styled" className="custom-class" />);
    
    const buttonElement = screen.getByTestId('custom-button');
    expect(buttonElement.className).toContain('custom-class');
  });

  // Example of a snapshot test
  test('matches snapshot', () => {
    const { container } = render(<Button text="Snapshot" />);
    expect(container).toMatchSnapshot();
  });
});

// Example of testing async behavior
describe('Button Component Async Behavior', () => {
  test('shows loading state when clicked', async () => {
    // Mock an async function
    const asyncHandleClick = jest.fn().mockImplementation(() => {
      return new Promise(resolve => {
        setTimeout(resolve, 100);
      });
    });
    
    const { rerender } = render(
      <Button text="Submit" onClick={asyncHandleClick} />
    );
    
    const buttonElement = screen.getByTestId('custom-button');
    userEvent.click(buttonElement);
    
    expect(asyncHandleClick).toHaveBeenCalledTimes(1);
    
    // Example of updating props to show loading state
    // (would require Button to accept an isLoading prop)
    rerender(<Button text="Loading..." disabled={true} />);
    
    expect(buttonElement).toHaveTextContent('Loading...');
    expect(buttonElement).toBeDisabled();
    
    // Wait for async operation to complete
    await waitFor(() => {
      expect(asyncHandleClick).toHaveBeenCalledTimes(1);
    });
  });
});