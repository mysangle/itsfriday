import { supabaseClient } from "@/store";
import { toItsError } from "./itserror";

export interface LibroYearMonthReport {
  yearMonth: string;
  count: number;
}

export interface MoneroYearMonthCategoryReport {
  yearMonth: string;
  categoryId: number;
  price: number;
}

export interface Exchange {
  id?: number;
  from: string;
  to: string;
  value: string;
  createdAt?: Date;
  updatedAt?: Date;
}

const homeStore = (() => {
  const fetchExchanges = async () => {
    let { data, error } = await supabaseClient
      .from("exchange")
      .select('*')
    return { data, error: error != null ? toItsError(error) : null };
  };

  const fetchBooksReadCountThisYear = async () => {
    let { data, error } = await supabaseClient
        .rpc('get_this_year_read_count')
    return { data, error: error != null ? toItsError(error) : null }
  };

  const fetchBooksReadCountByMonth = async () => {
    let { data, error } = await supabaseClient
        .rpc('get_year_month_count_stats')
    return { data, error: error != null ? toItsError(error) : null }
  };

  const fetchTotalExpensesThisYear = async () => {
    let { data, error } = await supabaseClient
        .rpc('stats_expense_this_year_total_price')
    return { data, error: error != null ? toItsError(error) : null }
  }

  const fetchPricesByMonthAndCategory = async () => {
    let { data, error } = await supabaseClient
        .rpc('stats_expense_year_month_price')
    return { data, error: error != null ? toItsError(error) : null }
  };

  return {
    fetchExchanges,
    fetchBooksReadCountThisYear,
    fetchBooksReadCountByMonth,
    fetchTotalExpensesThisYear,
    fetchPricesByMonthAndCategory,
  };
})();

export { homeStore }
