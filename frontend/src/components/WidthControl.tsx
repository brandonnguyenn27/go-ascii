import { useState, useEffect } from 'react';
import { Slider } from './ui/slider';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { Switch } from './ui/switch';

interface WidthControlProps {
  value: number;
  onChange: (value: number) => void;
  enabled: boolean;
  onToggle: (enabled: boolean) => void;
  min?: number;
  max?: number;
}

export function WidthControl({
  value,
  onChange,
  enabled,
  onToggle,
  min = 20,
  max = 300,
}: WidthControlProps) {
  const [inputValue, setInputValue] = useState(value.toString());

  useEffect(() => {
    setInputValue(value.toString());
  }, [value]);

  const handleSliderChange = (newValue: number[]) => {
    onChange(newValue[0]);
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);
  };

  const handleInputBlur = () => {
    const numValue = parseInt(inputValue, 10);
    if (!isNaN(numValue)) {
      const clampedValue = Math.max(min, Math.min(max, numValue));
      onChange(clampedValue);
    } else {
      setInputValue(value.toString());
    }
  };

  const handleInputKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleInputBlur();
    }
  };

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <Switch
            id="width-toggle"
            checked={enabled}
            onCheckedChange={onToggle}
          />
          <Label htmlFor="width-toggle" className="cursor-pointer">
            Custom Width
          </Label>
        </div>
        {enabled && (
          <Input
            id="width-control"
            type="number"
            min={min}
            max={max}
            value={inputValue}
            onChange={handleInputChange}
            onBlur={handleInputBlur}
            onKeyDown={handleInputKeyDown}
            className="w-20"
          />
        )}
      </div>
      {enabled && (
        <>
          <Slider
            value={[value]}
            onValueChange={handleSliderChange}
            min={min}
            max={max}
            step={1}
            className="w-full"
          />
          <p className="text-xs text-muted-foreground">
            Range: {min} - {max} characters
          </p>
        </>
      )}
      {!enabled && (
        <p className="text-xs text-muted-foreground">
          Using default width (100 characters)
        </p>
      )}
    </div>
  );
}

