import { observer } from "mobx-react-lite";
import { useEffect, useState } from "react";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Legend, ResponsiveContainer, Tooltip } from 'recharts';
import { toast } from "react-hot-toast";
import { Separator } from "@/components/ui/separator";
import { type LibroYearMonthReport, homeStore } from "@/types/model/home_service";
import { toCamelCase } from "@/utils/common";
import { useTranslate } from "@/utils/i18n";

// https://m.search.naver.com/p/csearch/content/qapirender.nhn\?key\=calculator\&pkid\=141\&q\=%ED%99%98%EC%9C%A8\&where\=m\&u1\=keb\&u6\=standardUnit\&u7\=0\&u3\=USD\&u4\=KRW\&u8\=down\&u2\=1

const LibroSection = observer(() => {
  const t = useTranslate();
  const [count, setCount] = useState<number>(0);
  const [reports, setReports] = useState<LibroYearMonthReport[]>([]);

  useEffect(() => {
    fetchThisYearCount();
  }, []);

  useEffect(() => {
    fetchReports();
  }, []);

  function fillEmptyMonths(data: any): LibroYearMonthReport[] {
    // 1. create a list of months for the past 12 months
    const months: string[] = [];
    const today = new Date();
    for (let i = -1; i < 12; i++) {
      const d = new Date(today.getFullYear(), today.getMonth() - i, 1);
      const monthStr = d.toISOString().slice(0, 7); // "YYYY-MM"
      months.push(monthStr);
    }

    // 2. convert existing data into a Map
    const dataMap = new Map((data as LibroYearMonthReport[]).map(d => [d.yearMonth, d]));

    // 3. create the final array filled with zeros
    const filledData: LibroYearMonthReport[] = months.map(m => ({
      yearMonth: m,
      count: dataMap.get(m)?.count ?? 0, // 0 if empty
    }));

    return filledData;
  }

  const fetchThisYearCount = async () => {
    try {
      let { data, error } = await homeStore.fetchBooksReadCountThisYear()
      if (error != null) {
        throw error;
      }
      setCount(data);
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
    }
  };

  const fetchReports = async () => {
    try {
      let { data, error } = await homeStore.fetchBooksReadCountByMonth()
      if (error != null) {
        throw error;
      }

      if (data != null) {
        const filledData = fillEmptyMonths((data as LibroYearMonthReport[]).map((r) => toCamelCase(r)));
        setReports(filledData);
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
    }
  };

  return (
    <div className="w-144 h-80 flex flex-col border border-dashed rounded-md justify-center items-start px-4 py-2 text-muted-foreground">
      <div className="flex items-center justify-between">
        <div className="flex-auto space-y-1">
          <p className="flex flex-row justify-start items-center font-medium text-muted-foreground">
            Libro Report
          </p>
          <p className="text-sm text-muted-foreground">
            {t("libro.books-read-count-this-year")}: {count}
          </p>
        </div>
      </div>
      <Separator className="mt-2 mb-4" />
      <ResponsiveContainer width="100%" height="70%">
        <BarChart
          data={reports}
          margin={{
            top: 5,
            right: 30,
            left: 20,
            bottom: 5,
          }}
        >
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="yearMonth" />
          <YAxis />
          <Tooltip cursor={false} />
          <Legend />
          <Bar dataKey="count" name={t("libro.monthly-books-read-count")} fill="#82ca9d" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
});

export default LibroSection;
