import { MoreVerticalIcon, PlusIcon } from "lucide-react";
import { observer } from "mobx-react-lite";
import { useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import { Button } from "@/components/ui/button";
import { useDialog } from "@/hooks/useDialog";
import CreateCategoryDialog from "../CreateCategoryDialog";
import { type ExpenseCategory, dineroStore } from "@/types/model/dinero_service";
import { toCamelCase } from "@/utils/common";
import { useTranslate } from "@/utils/i18n";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu";

const CategorySection = observer(() => {
  const t = useTranslate();
  const [categories, setCategories] = useState<ExpenseCategory[]>([]);
  const createDialog = useDialog();
  const editDialog = useDialog();
  const [editingCategory, setEditingCategory] = useState<ExpenseCategory | undefined>();

  useEffect(() => {
    fetchCategories();
  }, []);

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

  const handleCreateCategory = () => {
    setEditingCategory(undefined);
    createDialog.open();
  };

  const handleEditCategory = (category: ExpenseCategory) => {
    setEditingCategory(category);
    editDialog.open();
  };

  const handleDeleteCategory = async (category: ExpenseCategory) => {
    const confirmed = window.confirm(t("common.delete-warning", { title: category.name }));
    if (confirmed) {
      if (category.id !== undefined) {
        let { error } = await dineroStore.deleteCategory(category.id)
        if (error != null) {
          console.error(error);
        }
        fetchCategories();
      }
    }
  };

  return (
    <div className="w-full flex flex-col gap-2 pt-2 pb-4">
      <div className="w-full flex flex-col flex-row gap-2 pt-4 pb-4 justify-between items-center">
        <p className="font-medium text-muted-foreground">{t("dinero.create-a-category")}</p>
        <Button onClick={handleCreateCategory}>
          <PlusIcon className="w-4 h-4 mr-2" />
          {t("common.create")}
        </Button>
      </div>
      <div className="w-full flex flex-row justify-between items-center mt-6">
        <div className="title-text">{t("dinero.category-list")}</div>
      </div>
      <div className="w-full overflow-x-auto">
        <div className="inline-block min-w-full align-middle border border-border rounded-lg">
          <table className="min-w-full divide-y divide-border">
            <thead>
              <tr className="text-sm font-semibold text-left text-foreground">
                <th scope="col" className="px-3 py-2 text-right">
                  {t("dinero.id")}
                </th>
                <th scope="col" className="px-3 py-2">
                  {t("dinero.name")}
                </th>
                <th scope="col" className="relative py-2 pl-3 pr-4"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {categories.map((category) => (
                <tr key={category.id} className="text-left">
                  <td className="px-3 py-2 text-right text-sm text-muted-foreground">{category.id}</td>
                  <td className="px-3 py-2 text-sm text-muted-foreground">{category.name}</td>
                  <td className="relative whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium flex justify-end">
                    <DropdownMenu modal={false}>
                      <DropdownMenuTrigger asChild>
                        <Button variant="outline">
                          <MoreVerticalIcon className="w-4 h-auto" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end" sideOffset={2}>
                        <>
                          <DropdownMenuItem onClick={() => handleEditCategory(category)}>{t("common.update")}</DropdownMenuItem>
                          <DropdownMenuItem
                            onClick={() => handleDeleteCategory(category)}
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

      {/* Create Category Dialog */}
      <CreateCategoryDialog open={createDialog.isOpen} onOpenChange={createDialog.setOpen} onSuccess={fetchCategories} />

      {/* Edit Category Dialog */}
      <CreateCategoryDialog open={editDialog.isOpen} onOpenChange={editDialog.setOpen} category={editingCategory} onSuccess={fetchCategories} />
    </div>
  );
});

export default CategorySection;
