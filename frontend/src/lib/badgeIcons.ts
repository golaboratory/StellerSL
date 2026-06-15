// Maps a badge requirement_type to its bundled SVG asset.
// SVG filenames in src/assets/badges/ match requirement_type exactly.
const svgs = import.meta.glob('../assets/badges/*.svg', {
  eager: true,
  query: '?url',
  import: 'default',
}) as Record<string, string>;

export function badgeIconUrl(requirementType?: string): string | null {
  if (!requirementType) return null;
  return svgs[`../assets/badges/${requirementType}.svg`] ?? null;
}
