import dayjs from "dayjs";
import { observer } from "mobx-react-lite";
import { useState } from "react";
import { cn } from "@/lib/utils";
import ActivityCalendar from "../components/ActivityCalendar";
import { MonthNavigator } from "../components/StatisticsView/MonthNavigator";

const Home = observer(() => {
  const [selectedDate] = useState(new Date());
  const [visibleMonthString, setVisibleMonthString] = useState(dayjs().format("YYYY-MM"));

  return (
    <div className="w-full min-h-full bg-background text-foreground">
      <div className="text-3xl font-bold underline">It's Friday</div>
      <div className={cn("shrink-0 h-svh transition-all", "w-72", "group mt-2 px-3 py-6  space-y-1 text-muted-foreground animate-fade-in")}>
        <MonthNavigator visibleMonth={visibleMonthString} onMonthChange={setVisibleMonthString} />
        <div className="animate-scale-in">
          <ActivityCalendar
            month={visibleMonthString}
            selectedDate={selectedDate.toDateString()}
            data={{}}
          /> 
        </div>
      </div>
    </div>
  );
});

export default Home;
