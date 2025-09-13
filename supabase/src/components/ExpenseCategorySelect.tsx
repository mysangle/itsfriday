import { GlobeIcon } from "lucide-react";
import { type FC } from "react";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { type ExpenseCategory } from "@/types/model/dinero_service";

interface Props {
  categories: ExpenseCategory[],
  value: ExpenseCategory;
  onChange: (category: ExpenseCategory) => void;
}

const ExpenseCategorySelect: FC<Props> = (props: Props) => {
  const { categories, onChange, value } = props;

  const handleSelectChange = async (name: string) => {
    const category = categories.find(c => c.name === name)
    if (category !== undefined) {
      onChange(category);
    }
  };

  return (
    <Select value={value.name} onValueChange={handleSelectChange}>
      <SelectTrigger>
        <div className="flex items-center gap-2">
          <GlobeIcon className="w-4 h-auto" />
          <SelectValue placeholder="Select category" />
        </div>
      </SelectTrigger>
      <SelectContent>
        {categories.map((category) => {
          return (
            <SelectItem key={category.name} value={category.name}>
              {category.name}
            </SelectItem>
          );
        })}
      </SelectContent>
    </Select>
  );
};

export default ExpenseCategorySelect;
