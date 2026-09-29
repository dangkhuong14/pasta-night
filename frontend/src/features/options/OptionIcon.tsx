import {
  Clapperboard,
  Film,
  Heart,
  Popcorn,
  User,
  Users,
  type LucideProps,
} from "lucide-react";

type OptionIconProps = LucideProps & {
  /** The API's option `icon`: a lucide icon name (ARCHITECTURE §1). */
  name: string;
};

/**
 * Renders the lucide icon for an option, `Film` when the name is unknown.
 * A static switch keeps icons server-rendered with no client JS; add a case
 * when options.yaml gains a new icon name.
 */
export function OptionIcon({ name, ...props }: OptionIconProps) {
  switch (name) {
    case "heart":
      return <Heart {...props} />;
    case "user":
      return <User {...props} />;
    case "users":
      return <Users {...props} />;
    case "popcorn":
      return <Popcorn {...props} />;
    case "clapperboard":
      return <Clapperboard {...props} />;
    default:
      return <Film {...props} />;
  }
}
