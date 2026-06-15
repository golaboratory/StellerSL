// Lightweight form validation helpers (no extra dependencies).

export function isValidEmail(value: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value);
}

// Returns an error message, or empty string when the password is acceptable:
// at least 8 characters, containing both a letter and a digit.
export function passwordStrengthError(value: string): string {
  if (value.length < 8) return 'Password must be at least 8 characters';
  if (!/[A-Za-z]/.test(value) || !/[0-9]/.test(value)) {
    return 'Password must contain both letters and numbers';
  }
  return '';
}
