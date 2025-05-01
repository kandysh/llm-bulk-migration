import React from 'react';
import { mount, shallow } from 'enzyme';
import Button from '../component1';
import {capitalize, chunkArray} from "../utils/helper";

describe('Button Component (Enzyme Tests)', () => {
  // Test that the component renders with the correct text
  it('renders button with correct text', () => {
    const wrapper = shallow(<Button text="Click me" />);
    expect(wrapper.text()).toContain('Click me');
    expect(wrapper.find('button').prop('data-testid')).toBe('custom-button');
  });

  // Test different variants of the button
  it('renders with primary variant by default', () => {
    const wrapper = shallow(<Button text="Primary" />);
    expect(wrapper.find('button').prop('className')).toContain('bg-blue-500');
  });

  it('renders with secondary variant when specified', () => {
    const wrapper = shallow(<Button text="Secondary" variant="secondary" />);
    expect(wrapper.find('button').prop('className')).toContain('bg-gray-200');
  });

  it('renders with danger variant when specified', () => {
    const wrapper = shallow(<Button text="Danger" variant="danger" />);
    expect(wrapper.find('button').prop('className')).toContain('bg-red-500');
  });

  // Test the onClick handler
  it('calls onClick handler when clicked', () => {
    const handleClick = jest.fn();
    const wrapper = shallow(<Button text="Click me" onClick={handleClick} />);
    
    wrapper.find('button').simulate('click');
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  // Test disabled state
  it('does not call onClick when disabled', () => {
    const handleClick = jest.fn();
    const wrapper = shallow(
      <Button text="Click me" onClick={handleClick} disabled={true} />
    );
    
    wrapper.find('button').simulate('click');
    expect(handleClick).not.toHaveBeenCalled();
    expect(wrapper.find('button').prop('disabled')).toBe(true);
    expect(wrapper.find('button').prop('className')).toContain('opacity-50');
  });

  // Test hover state
  it('shows arrow on hover', () => {
    const wrapper = mount(<Button text="Hover me" />);
    
    // Initially, arrow should not be visible
    expect(wrapper.text()).not.toContain('→');
    
    // Simulate mouse enter
    wrapper.find('button').simulate('mouseenter');
    wrapper.update();
    
    // Arrow should now be visible
    expect(wrapper.text()).toContain('→');
    
    // Simulate mouse leave
    wrapper.find('button').simulate('mouseleave');
    wrapper.update();
    
    // Arrow should be hidden again
    expect(wrapper.text()).not.toContain('→');
  });

  // Test custom class names
  it('applies custom className', () => {
    const wrapper = shallow(<Button text="Styled" className="custom-class" />);
    expect(wrapper.find('button').prop('className')).toContain('custom-class');
  });

  // Test that useEffect is called on mount and unmount
  it('calls useEffect on mount and unmount', () => {
    // Spy on console.log
    const consoleSpy = jest.spyOn(console, 'log');
    
    const wrapper = mount(<Button text="Effect Test" />);
    expect(consoleSpy).toHaveBeenCalledWith('Button component mounted');
    
    // Clean up
    wrapper.unmount();
    expect(consoleSpy).toHaveBeenCalledWith('Button component will unmount');
    
    consoleSpy.mockRestore();
  });

  // Snapshot test
  it('matches snapshot', () => {
    const wrapper = shallow(<Button text="Snapshot" />);
    expect(wrapper).toMatchSnapshot();
  });
});

// Mock the useState hook for more complex tests
describe('Button Component State Tests', () => {
  it('sets isHovered state correctly', () => {
    // Create a mock for useState
    const mockSetState = jest.fn();
    const useStateMock = jest.spyOn(React, 'useState');
    
    // Mock the first call to useState for isHovered
    useStateMock.mockImplementationOnce(() => [false, mockSetState]);
    
    const wrapper = shallow(<Button text="State Test" />);
    
    // Simulate mouseenter and check if setState was called with true
    wrapper.find('button').simulate('mouseenter');
    expect(mockSetState).toHaveBeenCalledWith(true);
    
    // Simulate mouseleave and check if setState was called with false
    wrapper.find('button').simulate('mouseleave');
    expect(mockSetState).toHaveBeenCalledWith(false);
    
    // Clean up
    useStateMock.mockRestore();
  });
});