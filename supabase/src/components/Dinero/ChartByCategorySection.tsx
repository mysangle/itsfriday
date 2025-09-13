import { observer } from "mobx-react-lite";
import { useEffect, useState } from "react";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Legend, ResponsiveContainer, Tooltip } from 'recharts';
import { toast } from "react-hot-toast";
import { Separator } from "@/components/ui/separator";
import { type DineroYearMonthCategoryReport, homeStore } from "@/types/model/home_service";
import { type ExpenseCategory, dineroStore } from "@/types/model/dinero_service";
import { toCamelCase } from "@/utils/common";
import { useTranslate } from "@/utils/i18n";

type ReportItem = Map<string, any>
const COLORS = ['#0088FE', '#ff5242ff', '#FFBB28', '#8884d8', '#00C49F'];

const ChartByCategorySection = observer(() => {
  const t = useTranslate();
  const [reports, setReports] = useState<Object[]>([]);
  const [categories, setCategories] = useState<ExpenseCategory[]>([]);

  useEffect(() => {
    fetchCategories();
  }, []);

  useEffect(() => {
    if (categories.length > 0) {
      fetchReports();
    }
  }, [categories]);

  const fetchCategories = async () => {
    try {
      let { data, error } = await dineroStore.fetchCategories()
      if (error != null) {
        throw error;
      }

      if (data) {
        setCategories(data.map((category) => (toCamelCase(category))))
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
    }
  };

  function convertToReportData(items: DineroYearMonthCategoryReport[]): ReportItem {
    const data = new Map<string, any>();
    for (const category of categories) {
      if (category.id !== undefined) {
        data.set(category.name, 0);
        for (const item of items) {
          if (category.id === item.categoryId) {
            data.set(category.name, item.price);
          }
        }
      }
    }
    return data;
  }

  function fillEmptyMonths(data: DineroYearMonthCategoryReport[]): Object[] {
    // 1. create a list of months for the past 12 months
    const months: string[] = [];
    const monthSet = new Map<string, DineroYearMonthCategoryReport[]>();
    const today = new Date();
    for (let i = -1; i < 12; i++) {
      const d = new Date(today.getFullYear(), today.getMonth() - i, 1);
      const monthStr = d.toISOString().slice(0, 7); // "YYYY-MM"
      months.push(monthStr);

      const dataArray = data.filter(d => d.yearMonth === monthStr)
      monthSet.set(monthStr, dataArray);
    }

    const filledData: ReportItem[] = []
    for (const month of months) {
      const dataArray = monthSet.get(month);
      if (dataArray !== undefined) {
        const data = convertToReportData(dataArray);
        data.set("yearMonth", month);
        filledData.push(data);
      }
    }

    return filledData.map(m => Object.fromEntries(m));
  }

  const fetchReports = async () => {
    try {
      let { data, error } = await homeStore.fetchPricesByMonthAndCategory()
      if (error != null) {
        throw error;
      }

      if (data != null && categories.length > 0) {
        const filledData = fillEmptyMonths((data as DineroYearMonthCategoryReport[]).map((r) => toCamelCase(r)));
        setReports(filledData);
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
    }
  };

  return (
    <div className="w-full h-full flex flex-col rounded-md justify-center items-start px-4 py-2 text-muted-foreground">
      <div className="flex items-center justify-start">
        <div className="space-y-1">
          <p className="flex flex-row justify-start items-center font-medium text-muted-foreground">
            By Category
          </p>
          <p className="text-sm text-left text-muted-foreground">
            {t("dinero.monthly-spending-by-category")}
          </p>
        </div>
      </div>
      <Separator className="mt-2 mb-4" />
      <div className="w-full h-120 border border-dashed rounded-md px-4 py-2">
        <ResponsiveContainer width="100%" height="100%">
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
            {categories.map((entry, index) => (
              <Bar key={index} dataKey={entry.name} name={entry.name} fill={COLORS[index % COLORS.length]} />
            ))}
          </BarChart>
        </ResponsiveContainer>
      </div>

    </div>
  );
});

export default ChartByCategorySection;
