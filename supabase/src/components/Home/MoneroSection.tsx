import { observer } from "mobx-react-lite";
import { useEffect, useState } from "react";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Legend, ResponsiveContainer, Tooltip } from 'recharts';
import { toast } from "react-hot-toast";
import { Separator } from "@/components/ui/separator";
import { thisMonth } from "@/helpers/utils";
import { type MoneroYearMonthCategoryReport, homeStore } from "@/types/model/home_service";
import { type ExpenseCategory, moneroStore } from "@/types/model/monero_service";
import { toCamelCase } from "@/utils/common";
import { useTranslate } from "@/utils/i18n";

type ReportItem = Map<string, any>
const COLORS = ['#0088FE', '#ff5242ff', '#FFBB28', '#8884d8', '#00C49F'];

const MoneroSection = observer(() => {
  const t = useTranslate();
  const [totalPrice, setTotalPrice] = useState<number>(0);
  const [totalPriceThisMonth, setTotalPriceThisMonth] = useState<number>(0);
  const [reports, setReports] = useState<Object[]>([]);
  const [categories, setCategories] = useState<ExpenseCategory[]>([]);

  useEffect(() => {
    fetchThisYearTotalExpenses();
  }, []);

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
      let { data, error } = await moneroStore.fetchCategories()
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

  function convertToReportData(items: MoneroYearMonthCategoryReport[]): ReportItem {
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

  function calculateTotalSpendingThisMonth(items: MoneroYearMonthCategoryReport[]) {
    let price = 0;
    for (const item of items) {
      price += item.price;
    }
    setTotalPriceThisMonth(price);
  }

  function fillEmptyMonths(data: MoneroYearMonthCategoryReport[]): Object[] {
    // 1. create a list of months for the past 12 months
    const months: string[] = [];
    const monthSet = new Map<string, MoneroYearMonthCategoryReport[]>();
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

      if (month === thisMonth) {
        calculateTotalSpendingThisMonth(data);
      }
    }

    return filledData.map(m => Object.fromEntries(m));
  }

  const fetchThisYearTotalExpenses = async () => {
    try {
      let { data, error } = await homeStore.fetchTotalExpensesThisYear()
      if (error != null) {
        throw error;
      }
      setTotalPrice(data);
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
    }
  };

  const fetchReports = async () => {
    try {
      let { data, error } = await homeStore.fetchPricesByMonthAndCategory()
      if (error != null) {
        throw error;
      }

      if (data != null && categories.length > 0) {
        const filledData = fillEmptyMonths((data as MoneroYearMonthCategoryReport[]).map((r) => toCamelCase(r)));
        setReports(filledData);
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
    }
  };

  return (
    <div className="w-144 h-80 flex flex-col border border-dashed rounded-md justify-center items-start px-4 py-2 text-muted-foreground">
      <div className="flex items-center justify-start">
        <div className="space-y-1">
          <p className="flex flex-row justify-start items-center font-medium text-muted-foreground">
            Monero Report
          </p>
          <p className="text-sm text-left text-muted-foreground">
            {t("monero.total-price-this-year")}: {totalPrice.toLocaleString()}
          </p>
          <p className="text-sm text-left text-muted-foreground">
            {t("monero.total-price-this-month")}: {totalPriceThisMonth.toLocaleString()}
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
          {categories.map((entry, index) => (
            <Bar key={index} dataKey={entry.name} stackId="b" name={entry.name} fill={COLORS[index % COLORS.length]} />
          ))}
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
});

export default MoneroSection;
