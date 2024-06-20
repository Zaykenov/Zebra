export function getDateAndTime(dateString: string) {
  const date = new Date(dateString);
  const dateStr = date.toLocaleDateString();
  const timeStr = date.toLocaleTimeString();
  return `${dateStr} ${timeStr}`;
}

export function getDateString() {
  const today = new Date();
  const year = today.getFullYear();
  const month = String(today.getMonth() + 1).padStart(2, "0");
  const day = String(today.getDate()).padStart(2, "0");

  return `${year}-${month}-${day}`;
}

export const formatSeconds = (seconds: number) => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
};
