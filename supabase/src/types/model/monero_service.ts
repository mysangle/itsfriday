import { supabaseClient } from "@/store";
import { toItsError } from "./itserror";

export interface ExpenseRow {
  item: string;
  price: number;
  categoryId: number;
  dateUsed: string;
}

export function createExpenseRow(expense: Expense): ExpenseRow {
  if (!expense.price || !expense.dateUsed || !expense.expenseCategory || !expense.expenseCategory.id) {
    throw "invalid Expense";
  }

  return {
    item: expense.item,
    price: expense.price,
    categoryId: expense.expenseCategory.id,
    dateUsed: expense.dateUsed,
  };
}

export interface Expense {
  id?: number;
  item: string;
  price?: number;
  expenseCategory: ExpenseCategory;
  dateUsed?: string | undefined;
  createdAt?: Date | undefined;
}

function createBaseExpense(): Expense {
  return {
    item: "",
    expenseCategory: createBaseCategory(),
  };
}

export const Expense: MessageFns<Expense> = {
  create(base?: Partial<Expense>): Expense {
    return Expense.fromPartial(base ?? {});
  },
  fromPartial(object: Partial<Expense>): Expense {
    const message = createBaseExpense();
    message.id = object.id ?? undefined;
    message.item = object.item ?? "";
    message.price = object.price;
    if (object.expenseCategory) {
       message.expenseCategory = object.expenseCategory;
    }
    message.dateUsed = object.dateUsed ?? "";
    return message;
  },
};

export interface ExpenseCategory {
  id?: number;
  name: string;
  createdAt?: Date | undefined;
}

function createBaseCategory(): ExpenseCategory {
  return {
    name: "",
  };
}

export const ExpenseCategory: MessageFns<ExpenseCategory> = {
  create(base?: Partial<ExpenseCategory>): ExpenseCategory {
    return ExpenseCategory.fromPartial(base ?? {});
  },
  fromPartial(object: Partial<ExpenseCategory>): ExpenseCategory {
    const message = createBaseCategory();
    message.id = object.id ?? undefined;
    message.name = object.name ?? "";
    return message;
  },
};

export interface MessageFns<T> {
  create(base?: Partial<T>): T;
  fromPartial(object: Partial<T>): T;
}

const moneroStore = (() => {
  const insertCategory = async (category: any) => {
    const { error } = await supabaseClient
      .from("expense_category")
      .insert([category])
    return { error: error != null ? toItsError(error) : null };
  };

  const updateCategory = async (id: number, updatedCategory: any) => {
    const { error } = await supabaseClient
      .from("expense_category")
      .update(updatedCategory)
      .eq('id', id)
    return { error: error != null ? toItsError(error) : null };
  }

  const fetchCategories = async () => {
    let { data, error } = await supabaseClient
      .from("expense_category")
      .select('*')
      .order('created_at', { ascending: true });
    return { data, error: error != null ? toItsError(error) : null };
  };

  const deleteCategory = async (id: number) => {
    let { error } = await supabaseClient
      .from("expense_category")
      .delete()
      .eq('id', id);
    return { error: error != null ? toItsError(error) : null }
  };

  const fetchExpenses = async () => {
    let { data, error } = await supabaseClient
      .from("expense")
      .select(`
        id,
        item,
        price,
        date_used,
        created_at,
        expense_category (
          id,
          name
        )
      `)
      .order('date_used', { ascending: false });
    return { data, error: error != null ? toItsError(error) : null };
  };

  const deleteExpense = async (id: number) => {
    let { error } = await supabaseClient
      .from("expense")
      .delete()
      .eq('id', id);
    return { error: error != null ? toItsError(error) : null }
  };

  const insertExpense = async (expense: any) => {
    const { error } = await supabaseClient
      .from("expense")
      .insert([expense])
    return { error: error != null ? toItsError(error) : null };
  };

  const updateExpense = async (id: number, updatedExpense: any) => {
    const { error } = await supabaseClient
      .from("expense")
      .update(updatedExpense)
      .eq('id', id)
    return { error: error != null ? toItsError(error) : null };
  }

  return {
    insertCategory,
    updateCategory,
    fetchCategories,
    deleteCategory,
    fetchExpenses,
    deleteExpense,
    insertExpense,
    updateExpense,
  };
})();

export { moneroStore }
