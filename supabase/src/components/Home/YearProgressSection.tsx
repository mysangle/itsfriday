import { observer } from "mobx-react-lite";
import { useEffect, useState } from "react";
import { useTranslate } from "@/utils/i18n";

const YearProgressSection = observer(() => {
  const t = useTranslate();
  const [data, setData] = useState({ percentage: "0", elapsed: 0, remaining: 0, year: 2026 });

  useEffect(() => {
    const update = () => {
      const now = new Date();
      const year = now.getFullYear();
      const start = new Date(year, 0, 1);
      const end = new Date(year + 1, 0, 1);
      
      const total = end.getTime() - start.getTime();
      const elapsed = now.getTime() - start.getTime();
      const percentage = (elapsed / total) * 100;
      
      setData({
        percentage: percentage.toFixed(2),
        elapsed: Math.floor(elapsed / (1000 * 60 * 60 * 24)),
        remaining: Math.ceil((total - elapsed) / (1000 * 60 * 60 * 24)),
        year
      });
    };

    update();
  }, []);

  return (
    <div className="border border-dashed rounded-md flex flex-col justify-start items-center mt-4 px-4 py-2">
      <div className="text-sm font-small text-center mb-2">
        <p>{t("yearprogress.complete", { year: data.year, percentage: data.percentage})}</p>
        <div className="text-center text-sm">
          {t("yearprogress.elapsed", { elapsed: data.elapsed, percentage: data.percentage})}
        </div>
        <div className="text-center text-sm">
          {t("yearprogress.remaining", { remaining: data.remaining })}
        </div>
      </div>
      <div className="w-full bg-gray-200 rounded-sm h-8 mb-2">
        <div 
          className="bg-gradient-to-r from-purple-500 to-indigo-600 h-full rounded-sm transition-all duration-500"
          style={{ width: `${data.percentage}%` }}
        />
      </div>
    </div>
  );
});

export default YearProgressSection;
