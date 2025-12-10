import { Label } from './ui/label';

interface PaletteSelectorProps {
  value: string;
  onChange: (value: string) => void;
}

const PALETTE_OPTIONS = [
  { value: 'normal', label: 'Normal' },
  { value: 'dense', label: 'Dense' },
  { value: 'sparse', label: 'Sparse' },
  { value: 'unicode', label: 'Unicode' },
];

export function PaletteSelector({ value, onChange }: PaletteSelectorProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor="palette">Character Palette</Label>
      <select
        id="palette"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
      >
        {PALETTE_OPTIONS.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </div>
  );
}

