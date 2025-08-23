import { observer } from "mobx-react-lite";
import { useEffect, useState } from "react";
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts';
import { toast } from "react-hot-toast";
import { type CountByGenre, libroStore } from "@/types/model/libro_service";
import { useTranslate } from "@/utils/i18n";

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8'];

const ChartByGenreSection = observer(() => {
  const t = useTranslate();
  const [genreCount, setGenreCount] = useState<CountByGenre[]>([]);

  useEffect(() => {
    fetchStatsByGenre();
  }, []);

  const fetchStatsByGenre = async () => {
    try {
      let { data, error } = await libroStore.fetchStatsByGenre()
      if (error != null) {
        throw error;
      }

      if (data != null) {
        setGenreCount(data)
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
    }
  };

  return (
    <div className="w-full flex flex-col gap-2 pt-2 pb-4">
      <div className="w-full flex flex-col gap-2 pt-4 justify-start items-start">
        <p className="font-medium text-muted-foreground">{t("libro.by-genre")}</p>
        <p className="text-sm text-muted-foreground">
          {t("libro.books-read-count-by-genre")}
        </p>
      </div>
      <div className="w-full overflow-x-auto">
        <div className="inline-block min-w-full align-middle h-80">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart width={400} height={400}>
              <Pie
                nameKey="genre"
                dataKey="count"
                isAnimationActive={true}
                data={genreCount}
                cx="50%"
                cy="50%"
                outerRadius={128}
                fill="#82ca9d"
                label={({ name, value, percent }) => `${name}: ${value} (${((percent ?? 1) * 100).toFixed(0)}%)`}>
                {genreCount.map((entry, index) => (
                  <Cell key={`cell-${entry.genre}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
});

export default ChartByGenreSection;
