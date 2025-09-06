
export const today = new Date().toISOString().split("T")[0]; // YYYY-MM-DD
export const thisMonth = new Date().toISOString().slice(0, 7); // YYYY-MM

export function isValidDateFormat(dateStr: string): boolean {
    if (!/^\d{4}-\d{2}-\d{2}$/.test(dateStr)) {
      console.error(dateStr + "success")
      return false;
    }

    const [year, month, day] = dateStr.split("-").map(Number);
    const date = new Date(year, month - 1, day);
    return (
      date.getFullYear() === year &&
      date.getMonth() === month - 1 &&
      date.getDate() === day
    );
  }
