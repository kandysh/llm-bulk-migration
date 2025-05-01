import React from 'react';
import { mount, shallow } from 'enzyme';
import Button from '../component2';
import {capitalize, chunkArray} from "../utils/helper";

describe('Button Component (Enzyme Tests)', () => {
  // Test basic rendering
  it('renders without crashing', () => {
    const wrapper = shallow(<Button text="Button" />);
    expect(wrapper.exists()).toBe(true);
  });

  it('renders button with correct text', () => {
    const wrapper = shallow(<Button text="Hello World" />);
    expect(wrapper.text()).toBe('Hello World');
  });

  // Test props and their defaults
  it('uses default props when not provided', () => {
    const wrapper = shallow(<Button text="Default Props" />);
    const buttonProps = wrapper.find('button').props();
    
    expect(buttonProps.disabled).toBe(false);
    expect(buttonProps.className).toContain('bg-blue-500'); // primary variant
  });

  it('applies props correctly when provided', () => {
    const onClick = jest.fn();
    const wrapper = shallow(
      <Button 
        text="Props Test" 
        onClick={onClick} 
        variant="danger" 
        disabled={true} 
        className="test-class"
      />
    );
    
    const buttonProps = wrapper.find('button').props();
    
    expect(buttonProps.onClick).toBeDefined();
    expect(buttonProps.disabled).toBe(true);
    expect(buttonProps.className).toContain('bg-red-500'); // danger variant
    expect(buttonProps.className).toContain('test-class');
  });

  // Test event handlers
  it('handles click events when not disabled', () => {
    const onClick = jest.fn();
    const wrapper = shallow(<Button text="Click Test" onClick={onClick} />);
    
    wrapper.find('button').simulate('click');
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('does not trigger onClick when disabled', () => {
    const onClick = jest.fn();
    const wrapper = shallow(
      <Button text="Disabled Test" onClick={onClick} disabled={true} />
    );
    
    wrapper.find('button').simulate('click');
    expect(onClick).not.toHaveBeenCalled();
  });

  // Test state management
  it('changes hover state on mouse events', () => {
    const wrapper = mount(<Button text="Hover Test" />);
    
    // Initial state - no arrow
    expect(wrapper.find('span').exists()).toBe(false);
    
    // After mouse enter
    wrapper.find('button').simulate('mouseenter');
    expect(wrapper.find('span').exists()).toBe(true);
    expect(wrapper.find('span').text()).toBe('→');
    
    // After mouse leave
    wrapper.find('button').simulate('mouseleave');
    expect(wrapper.find('span').exists()).toBe(false);
  });

  // Test style variations
  it('applies correct styles for each variant', () => {
    // Primary variant
    let wrapper = shallow(<Button text="Primary" variant="primary" />);
    expect(wrapper.find('button').prop('className')).toContain('bg-blue-500');
    
    // Secondary variant
    wrapper = shallow(<Button text="Secondary" variant="secondary" />);
    expect(wrapper.find('button').prop('className')).toContain('bg-gray-200');
    
    // Danger variant
    wrapper = shallow(<Button text="Danger" variant="danger" />);
    expect(wrapper.find('button').prop('className')).toContain('bg-red-500');
  });

  it('applies disabled styles when disabled', () => {
    const wrapper = shallow(<Button text="Disabled" disabled={true} />);
    expect(wrapper.find('button').prop('className')).toContain('opacity-50');
    expect(wrapper.find('button').prop('className')).toContain('cursor-not-allowed');
  });

  // Test lifecycle methods
  it('logs to console on mount and unmount', () => {
    const consoleSpy = jest.spyOn(console, 'log');
    
    const wrapper = mount(<Button text="Lifecycle Test" />);
    expect(consoleSpy).toHaveBeenCalledWith('Button component mounted');
    
    wrapper.unmount();
    expect(consoleSpy).toHaveBeenCalledWith('Button component will unmount');
    
    consoleSpy.mockRestore();
  });

  // Snapshot testing
  it('matches snapshot', () => {
    const wrapper = shallow(<Button text="Snapshot Test" />);
    expect(wrapper).toMatchSnapshot();
  });
});