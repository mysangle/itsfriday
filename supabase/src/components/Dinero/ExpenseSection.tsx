import { MoreVerticalIcon, PlusIcon } from "lucide-react";
import { observer } from "mobx-react-lite";
import { useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import { Button } from "@/components/ui/button";
import { useDialog } from "@/hooks/useDialog";
import useNavigateTo from "@/hooks/useNavigateTo";
import { Routes } from "@/router";
import CreateExpenseDialog from "../CreateExpenseDialog";
import { type Expense, type ExpenseCategory, dineroStore } from "@/types/model/dinero_service";
import { toCamelCase } from "@/utils/common";
import { useTranslate } from "@/utils/i18n";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu";

const ReviewSection = observer(() => {
  const t = useTranslate();
  const navigateTo = useNavigateTo();
  const [categories, setCategories] = useState<ExpenseCategory[]>([]);
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const createDialog = useDialog();
  const editDialog = useDialog();
  const [editingExpense, setEditingExpense] = useState<Expense | undefined>();

  useEffect(() => {
    fetchCategories();
  }, []);

  useEffect(() => {
    if (categories.length > 0) {
      fetchExpenses();
    }
  }, [categories]);

  const fetchCategories = async () => {
    try {
      let { data, error } = await dineroStore.fetchCategories()
      if (error != null) {
        throw error;
      }

      if (data) {
        if (data.length === 0) {
          navigateTo(Routes.DINERO + "#category");
          return;
        }
        setCategories(data.map((category) => (toCamelCase(category))))
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
    }
  };

  const fetchExpenses = async () => {
    try {
      let { data, error } = await dineroStore.fetchExpenses()
      if (error != null) {
        throw error;
      }

      if (data) {
        setExpenses(data.map((expense) => (toCamelCase(expense))))
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
    }
  };

  const handleCreateExpense = () => {
    setEditingExpense(undefined);
    createDialog.open();
  };

  const handleEditExpense = (expense: Expense) => {
    setEditingExpense(expense);
    editDialog.open();
  };

  const handleDeleteExpense = async (expense: Expense) => {
    const confirmed = window.confirm(t("common.delete-warning", { title: expense.item }));
    if (confirmed) {
      if (expense.id !== undefined) {
        let { error } = await dineroStore.deleteExpense(expense.id)
        if (error != null) {
          console.error(error);
        }
        fetchExpenses();
      }
    }
  };

  return (
    <div className="w-full flex flex-col gap-2 pt-2 pb-4">
      <div className="w-full flex flex-col flex-row gap-2 pt-4 pb-4 justify-between items-center">
        <p className="font-medium text-muted-foreground">{t("dinero.create-an-expense")}</p>
        <Button onClick={handleCreateExpense}>
          <PlusIcon className="w-4 h-4 mr-2" />
          {t("common.create")}
        </Button>
      </div>
      <div className="w-full flex flex-row justify-between items-center mt-6">
        <div className="title-text">{t("dinero.expense-list")}</div>
      </div>
      <div className="w-full overflow-x-auto">
        <div className="inline-block min-w-full align-middle border border-border rounded-lg">
          <table className="min-w-full divide-y divide-border">
            <thead>
              <tr className="text-sm font-semibold text-left text-foreground">
                <th scope="col" className="px-3 py-2">
                  {t("dinero.category")}
                </th>
                <th scope="col" className="px-3 py-2">
                  {t("libro.date-read")}
                </th>
                <th scope="col" className="px-3 py-2">
                  {t("dinero.item")}
                </th>
                <th scope="col" className="px-3 py-2">
                  {t("dinero.price")}
                </th>
                <th scope="col" className="relative py-2 pl-3 pr-4"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {expenses.map((expense) => (
                <tr key={expense.id} className="text-left">
                  <td className="px-3 py-2 text-sm text-muted-foreground">{expense.expenseCategory.name}</td>
                  <td className="whitespace-nowrap px-3 py-2 text-sm text-muted-foreground">{expense.dateUsed}</td>
                  <td className="whitespace-nowrap px-3 py-2 text-sm text-muted-foreground">{expense.item}</td>
                  <td className="px-3 py-2 text-sm text-muted-foreground">{expense.price?.toLocaleString()}</td>
                  <td className="relative whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium flex justify-end">
                    <DropdownMenu modal={false}>
                      <DropdownMenuTrigger asChild>
                        <Button variant="outline">
                          <MoreVerticalIcon className="w-4 h-auto" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end" sideOffset={2}>
                        <>
                          <DropdownMenuItem onClick={() => handleEditExpense(expense)}>{t("common.update")}</DropdownMenuItem>
                          <DropdownMenuItem
                            onClick={() => handleDeleteExpense(expense)}
                            className="text-destructive focus:text-destructive"
                          >
                            {t("common.delete")}
                          </DropdownMenuItem>
                        </>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Create Expense Dialog */}
      <CreateExpenseDialog open={createDialog.isOpen} onOpenChange={createDialog.setOpen} expenseCategories={categories} onSuccess={fetchExpenses} />

      {/* Edit Expense Dialog */}
      <CreateExpenseDialog open={editDialog.isOpen} onOpenChange={editDialog.setOpen} expense={editingExpense} expenseCategories={categories} onSuccess={fetchExpenses} />
    </div>
  );
});

export default ReviewSection;
