import { useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import useLoading from "@/hooks/useLoading";
import { ExpenseCategory, dineroStore } from "@/types/model/dinero_service";
import { toSnakeCase } from "@/utils/common";
import { useTranslate } from "@/utils/i18n";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  category?: ExpenseCategory;
  onSuccess?: () => void;
}

function CreateCategoryDialog({ open, onOpenChange, category: initialCategory, onSuccess }: Props) {
  const t = useTranslate();
  const [category, setCategory] = useState(ExpenseCategory.fromPartial({ ...initialCategory }));
  const requestState = useLoading(false);
  const isCreating = !initialCategory;

  useEffect(() => {
    if (initialCategory) {
      setCategory(ExpenseCategory.fromPartial(initialCategory));
    } else {
      setCategory(ExpenseCategory.fromPartial({}));
    }
  }, [initialCategory]);

  const setPartialCategory = (state: Partial<ExpenseCategory>) => {
    setCategory({
      ...category,
      ...state,
    });
  };

  const handleConfirm = async () => {
    if (isCreating && (!category.name)) {
      toast.error("Name cannot be empty");
      return;
    }

    try {
      requestState.setLoading();
      if (isCreating) {
        let { error } = await dineroStore.insertCategory(toSnakeCase(category))
        if (error != null) {
          throw error;
        }
        toast.success("Create a category successfully");
      } else {
        const updateCategory: Record<string, any> = {};
        if (category.name !== initialCategory?.name) {
          updateCategory.name = category.name;
        }
        console.error(category.id);
        if (category.id !== undefined) {
          let { error } = await dineroStore.updateCategory(category.id, toSnakeCase(updateCategory))
          if (error != null) {
            throw error;
          }
          toast.success("Update the category successfully");
        }
      }
      requestState.setFinish();
      onSuccess?.();
      onOpenChange(false);
      setCategory(ExpenseCategory.fromPartial({}));
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
          <DialogTitle>{`${isCreating ? t("common.create") : t("common.edit")} ${t("dinero.category")}`}</DialogTitle>
          <DialogDescription>create or update a category</DialogDescription>
        </DialogHeader>
        <div className="flex flex-col gap-4">
          <div className="grid gap-2">
            <Label htmlFor="title">{t("dinero.name")}</Label>
            <Input
              id="name"
              type="text"
              placeholder={t("dinero.name")}
              value={category.name}
              onChange={(e) =>
                setPartialCategory({
                  name: e.target.value,
                })
              }
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

export default CreateCategoryDialog;
