import { Switch } from './ui/switch';
import { Label } from './ui/label';

interface ColorToggleProps {
  enabled: boolean;
  onToggle: (enabled: boolean) => void;
}

export function ColorToggle({ enabled, onToggle }: ColorToggleProps) {
  return (
    <div className="flex items-center space-x-2">
      <Switch
        id="color-mode"
        checked={enabled}
        onCheckedChange={onToggle}
      />
      <Label htmlFor="color-mode" className="cursor-pointer">
        Color Mode
      </Label>
    </div>
  );
}

