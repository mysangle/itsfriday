import { type Day, startOfWeek, endOfWeek, startOfMonth, endOfMonth } from "date-fns";
import { observer } from "mobx-react-lite";
import { useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import { Calendar } from "@/components/ui/calendar"
import { Separator } from "@/components/ui/separator";
import { formatDateKey } from "@/helpers/utils";
import { type ExpenseByMonth, moneroStore } from "@/types/model/monero_service";
import { toCamelCase } from "@/utils/common";
import { useTranslate } from "@/utils/i18n";

const CalendarSection = observer(() => {
  const t = useTranslate();
  const [month, setMonth] = useState(new Date());
  const [date, setDate] = useState<Date | undefined>();
  const [expenses, setExpenses] = useState(new Map<string, number>());

  useEffect(() => {
    fetchExpensesByMonth();
  }, [month]);

  function getCalendarRange(month: Date, weekStartsOn: Day = 0) {
    const monthStart = startOfMonth(month);
    const calendarStart = startOfWeek(monthStart, { weekStartsOn });

    const monthEnd = endOfMonth(month);
    const calendarEnd = endOfWeek(monthEnd, { weekStartsOn });

    return { calendarStart, calendarEnd };
  }

  const fetchExpensesByMonth = async () => {
    try {
      let { calendarStart, calendarEnd } = getCalendarRange(month);
      let { data, error } = await moneroStore.fetchExpensesByStartAndEndDay(formatDateKey(calendarStart), formatDateKey(calendarEnd))
      if (error != null) {
        throw error;
      }

      if (data) {
        const expenseMap = (data as ExpenseByMonth[]).map((expense) => (toCamelCase(expense))).reduce((map, expense) => {
          map.set(expense.dateUsed, expense.price);
          return map;
        }, new Map<string, number>());
        setExpenses(expenseMap)
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
    }
  };

  // 커스텀 Day 컴포넌트
  const CustomDay = ({ ...props }) => {
    const dayDate: Date = props.day.date;
    const dateKey = formatDateKey(dayDate);
    const price = expenses.get(dateKey);

    const isSelected = date && formatDateKey(date) === dateKey;
    const isCurrentMonth = dayDate.getMonth() === props.day.displayMonth.getMonth();

    return (
      <div
        className={`
          relative w-full h-full p-1 rounded-md cursor-pointer transition-all hover:bg-gray-100
          ${isSelected ? 'bg-blue-600 text-white hover:bg-blue-700' : ''}
          ${!isCurrentMonth ? 'opacity-50' : ''}
        `}
        {...props}
      >
        <div className="flex flex-col items-center justify-center h-full">
          <div className={`text-sm font-medium ${isCurrentMonth ? '' : 'opacity-50'} ${isSelected ? 'text-white' : ''}`}>
            {dayDate.getDate()}
          </div>
          <div className={`text-xs px-1 py-0.5 rounded mt-1 ${isSelected
            ? 'bg-white text-blue-600'
            : `border`
            }`}>
            {(price && price.toLocaleString()) ?? 0}
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="w-full h-full flex flex-col rounded-md justify-center items-start px-4 py-2 text-muted-foreground">
      <div className="flex items-center justify-start">
        <div className="space-y-1">
          <p className="flex flex-row justify-start items-center font-medium text-muted-foreground">
            By Calendar
          </p>
          <p className="text-sm text-left text-muted-foreground">
            {t("monero.months-expenses")}
          </p>
        </div>
      </div>
      <Separator className="mt-2 mb-4" />
      <div className="w-full h-full flex justify-center border border-dashed rounded-md px-4 py-2">
        <Calendar
          mode="single"
          selected={date}
          onSelect={setDate}
          month={month}
          onMonthChange={setMonth}
          showOutsideDays={true}
          className="w-128 rounded-md border shadow-sm"
          captionLayout="dropdown"
          components={{
            Day: CustomDay
          }}
        />
      </div>
    </div>
  );
});

export default CalendarSection;
