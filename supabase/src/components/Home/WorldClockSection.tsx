import { observer } from "mobx-react-lite";

type City = {
  name: string;
  timeZone: string;
};

const cities: City[] = [
  { name: "San Francisco", timeZone: "America/Los_Angeles" },
  { name: "London", timeZone: "Europe/London" },
  { name: "Seoul", timeZone: "Asia/Seoul" },
  { name: "Sydney", timeZone: "Australia/Sydney" },
];

const WorldClockSection = observer(() => {

  function getCurrentDateTime(timeZone: string): string {
    return new Date().toLocaleString("en-US", {
      timeZone,
      // year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      // second: "2-digit",
    });
  }

  return (
    <div className="border border-dashed rounded-md flex flex-col justify-start items-start mt-4 px-2 py-2">
      {cities.map((city) => (
        <div key={city.name} className="w-full text-sm flex flex-col flex-row gap-2 px-2 justify-between items-center">
          <p className="font-small text-muted-foreground">{city.name}</p>
          <div >{getCurrentDateTime(city.timeZone)}</div>
        </div>
      ))}
    </div>
  );
});

export default WorldClockSection;
