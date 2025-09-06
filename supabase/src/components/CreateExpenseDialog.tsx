import { useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import ExpenseCategorySelect from "@/components/ExpenseCategorySelect";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { isValidDateFormat, today } from "@/helpers/utils";
import useLoading from "@/hooks/useLoading";
import { Expense, type ExpenseCategory, createExpenseRow, moneroStore } from "@/types/model/monero_service";
import { toSnakeCase } from "@/utils/common";
import { useTranslate } from "@/utils/i18n";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  expense?: Expense;
  expenseCategories: ExpenseCategory[];
  onSuccess?: () => void;
}

function CreateExpenseDialog({ open, onOpenChange, expense: initialExpense, expenseCategories, onSuccess }: Props) {
  const t = useTranslate();
  const [expense, setExpense] = useState(Expense.fromPartial({ ...initialExpense }));
  const requestState = useLoading(false);
  const isCreating = !initialExpense;

  useEffect(() => {
    if (initialExpense) {
      setExpense(Expense.fromPartial(initialExpense));
    } else {
      setExpense(Expense.fromPartial({}));
    }
  }, [initialExpense]);

  const setPartialExpense = (state: Partial<Expense>) => {
    setExpense({
      ...expense,
      ...state,
    });
  };

  const handleCategorySelectChange = async (category: ExpenseCategory) => {
    setPartialExpense({
      expenseCategory: category,
    })
  };

  const onPriceInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const price = value === "" ? undefined : parseInt(value, 10) || undefined;

    setPartialExpense({
      price: price,
    })
  };

  const handleConfirm = async () => {
    if (isCreating && (!expense.item || !expense.price || !expense.dateUsed)) {
      toast.error("Item, price and date used cannot be empty");
      return;
    }
    if (expense.price && expense.price <= 0) {
      toast.error("Price should be more than 0");
      return;
    }
    if (expense.dateUsed !== undefined && !isValidDateFormat(expense.dateUsed)) {
      toast.error("Date used should be format like '" + today + "'");
      return;
    }

    try {
      requestState.setLoading();
      if (isCreating) {
        const { error } = await moneroStore.insertExpense(toSnakeCase(createExpenseRow(expense)))
        if (error != null) {
          throw error;
        }
        toast.success("Create expense successfully");
      } else {
        const updateExpense: Record<string, any> = {};
        if (expense.item !== initialExpense?.item) {
          updateExpense.title = expense.item;
        }
        if (expense.price !== initialExpense?.price) {
          updateExpense.price = expense.price;
        }
        if (expense.dateUsed !== initialExpense?.dateUsed) {
          updateExpense.dateUsed = expense.dateUsed;
        }
        if (expense.id !== undefined) {
          const { error } = await moneroStore.updateExpense(expense.id, toSnakeCase(updateExpense))
          if (error != null) {
            throw error;
          }
          toast.success("Update expense successfully");
        }
        
      }
      requestState.setFinish();
      onSuccess?.();
      onOpenChange(false);
      setExpense(Expense.fromPartial({}));
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
      requestState.setError();
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>{`${isCreating ? t("common.create") : t("common.edit")} ${t("monero.expense")}`}</DialogTitle>
          <DialogDescription>create or update a expense</DialogDescription>
        </DialogHeader>
        <div className="flex flex-col gap-4">
          <div className="grid gap-2">
            <Label htmlFor="title">{t("monero.category")}</Label>
            <ExpenseCategorySelect categories={expenseCategories} value={expense.expenseCategory} onChange={handleCategorySelectChange} />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="date_used">{t("monero.date-used") + "(format: " + today + ")"}</Label>
            <Input
              id="date_used"
              type="text"
              placeholder={today}
              value={expense.dateUsed}
              onChange={(e) =>
                setPartialExpense({
                  dateUsed: e.target.value,
                })
              }
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="title">{t("monero.item")}</Label>
            <Input
              id="item"
              type="text"
              placeholder={t("monero.item")}
              value={expense.item}
              onChange={(e) =>
                setPartialExpense({
                  item: e.target.value,
                })
              }
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="price">{t("monero.price")}</Label>
            <Input
              id="price"
              type="text"
              placeholder={"0"}
              min={0}
              value={expense.price ?? ''}
              onChange={onPriceInputChange}
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="ghost" disabled={requestState.isLoading} onClick={() => onOpenChange(false)}>
            {t("common.cancel")}
          </Button>
          <Button disabled={requestState.isLoading} onClick={handleConfirm}>
            {t("common.confirm")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

export default CreateExpenseDialog;
