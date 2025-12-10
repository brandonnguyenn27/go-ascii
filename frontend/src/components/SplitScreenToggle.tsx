import { Switch } from './ui/switch';
import { Label } from './ui/label';

interface SplitScreenToggleProps {
  enabled: boolean;
  onToggle: (enabled: boolean) => void;
}

export function SplitScreenToggle({ enabled, onToggle }: SplitScreenToggleProps) {
  return (
    <div className="flex items-center space-x-2">
      <Switch
        id="split-screen"
        checked={enabled}
        onCheckedChange={onToggle}
      />
      <Label htmlFor="split-screen" className="cursor-pointer">
        Split Screen
      </Label>
    </div>
  );
}

