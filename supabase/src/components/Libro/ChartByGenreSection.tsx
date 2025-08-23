import { observer } from "mobx-react-lite";
import { useEffect, useState } from "react";
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts';
import { toast } from "react-hot-toast";
import { supabaseClient } from "@/store";
import type { CountByGenre } from "@/types/model/libro_service";
import { useTranslate } from "@/utils/i18n";

type TooltipPayload = ReadonlyArray<any>;

type Coordinate = {
  x: number;
  y: number;
};

type PieSectorData = {
  percent?: number;
  name?: string | number;
  midAngle?: number;
  middleRadius?: number;
  tooltipPosition?: Coordinate;
  value?: number;
  paddingAngle?: number;
  dataKey?: string;
  payload?: any;
  tooltipPayload?: ReadonlyArray<TooltipPayload>;
};

type GeometrySector = {
  cx: number;
  cy: number;
  innerRadius: number;
  outerRadius: number;
  startAngle: number;
  endAngle: number;
};

type PieLabelProps = PieSectorData &
  GeometrySector & {
    tooltipPayload?: any;
  };

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8'];
const RADIAN = Math.PI / 180;
const renderCustomizedLabel = ({ cx, cy, midAngle, innerRadius, outerRadius, percent }: PieLabelProps) => {
  const radius = innerRadius + (outerRadius - innerRadius) * 0.5;
  const x = cx + radius * Math.cos(-(midAngle ?? 0) * RADIAN);
  const y = cy + radius * Math.sin(-(midAngle ?? 0) * RADIAN);

  return (
    <text x={x} y={y} fill="white" textAnchor={x > cx ? 'start' : 'end'} dominantBaseline="central">
      {`${((percent ?? 1) * 100).toFixed(0)}%`}
    </text>
  );
};

const ChartByGenreSection = observer(() => {
  const t = useTranslate();
  const [genreCount, setGenreCount] = useState<CountByGenre[]>([]);

  useEffect(() => {
    fetchStatsByGenre();
  }, []);

  const fetchStatsByGenre = async () => {
    try {
      let { data, error } = await supabaseClient
        .rpc('get_genre_12_stats')
      if (error != null) {
        throw error;
      }

      console.error(data)
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
          {t("libro.books-read-count-this-year-by-genre")}
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
