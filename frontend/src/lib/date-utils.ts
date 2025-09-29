import { format, isToday, isYesterday } from "date-fns"

export function formatRelativeDate(date: Date): string {
  if (isToday(date)) {
    return "today"
  } else if (isYesterday(date)) {
    return "yesterday"
  } else {
    // Format as "September 5th" with ordinal suffix
    const day = date.getDate()
    const month = format(date, "MMMM")
    
    // Add ordinal suffix (st, nd, rd, th)
    let suffix = "th"
    if (day < 11 || day > 13) {
      switch (day % 10) {
        case 1:
          suffix = "st"
          break
        case 2:
          suffix = "nd"
          break
        case 3:
          suffix = "rd"
          break
      }
    }
    
    return `${month} ${day}${suffix}`
  }
}