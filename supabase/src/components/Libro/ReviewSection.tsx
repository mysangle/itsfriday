import { MoreVerticalIcon, PlusIcon } from "lucide-react";
import { observer } from "mobx-react-lite";
import { useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import { Button } from "@/components/ui/button";
import { useDialog } from "@/hooks/useDialog";
import { type BookReview, libroStore } from "@/types/model/libro_service";
import { toCamelCase } from "@/utils/common";
import { useTranslate } from "@/utils/i18n";
import CreateReviewDialog from "../CreateReviewDialog";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu";

const ReviewSection = observer(() => {
  const t = useTranslate();
  const [reviews, setReviews] = useState<BookReview[]>([]);
  const createDialog = useDialog();
  const editDialog = useDialog();
  const [editingReview, setEditingReview] = useState<BookReview | undefined>();

  useEffect(() => {
    fetchReviews();
  }, []);

  const fetchReviews = async () => {
    try {
      let { data, error } = await libroStore.fetchBookReviews(50)
      if (error != null) {
        throw error;
      }

      if (data) {
        setReviews(data.map((review) => (toCamelCase(review))))
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
    }
  };

  const handleCreateReview = () => {
    setEditingReview(undefined);
    createDialog.open();
  };

  const handleEditReview = (review: BookReview) => {
    setEditingReview(review);
    editDialog.open();
  };

  const handleDeleteReview = async (review: BookReview) => {
    const confirmed = window.confirm(t("common.delete-warning", { title: review.title }));
    if (confirmed) {
      if (review.id !== undefined) {
        let { error } = await libroStore.deleteBookReview(review.id)
        if (error != null) {
          console.error(error);
        }
        fetchReviews();
      }
    }
  };

  return (
    <div className="w-full flex flex-col gap-2 pt-2 pb-4">
      <div className="w-full flex flex-col flex-row gap-2 pt-4 pb-4 justify-between items-center">
        <p className="font-medium text-muted-foreground">{t("libro.create-a-review")}</p>
        <Button onClick={handleCreateReview}>
          <PlusIcon className="w-4 h-4 mr-2" />
          {t("common.create")}
        </Button>
      </div>
      <div className="w-full flex flex-row justify-between items-center mt-6">
        <div className="title-text">{t("libro.review-list")}</div>
      </div>
      <div className="w-full overflow-x-auto">
        <div className="inline-block min-w-full align-middle border border-border rounded-lg">
          <table className="min-w-full divide-y divide-border">
            <thead>
              <tr className="text-sm font-semibold text-left text-foreground">
                <th scope="col" className="text-right px-3 py-2">
                  {t("libro.genre")}
                </th>
                <th scope="col" className="px-3 py-2">
                  {t("libro.title")}
                </th>
                <th scope="col" className="px-3 py-2">
                  {t("libro.author") + " / " + t("libro.translator")}
                </th>
                <th scope="col" className="px-3 py-2">
                  {t("libro.pub-year")}
                </th>
                <th scope="col" className="px-3 py-2">
                  {t("libro.date-read")}
                </th>
                <th scope="col" className="relative py-2 pl-3 pr-4"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {reviews.map((review) => (
                <tr key={review.id} className="text-left">
                  <td className="whitespace-nowrap px-3 py-2 text-right text-sm text-muted-foreground">{review.genre}</td>
                  <td className="px-3 py-2 text-sm text-muted-foreground">{review.title}</td>
                  <td className="px-3 py-2 text-sm text-muted-foreground">
                    {review.author + (review.translator ? " / " + review.translator : "")}
                  </td>
                  <td className="whitespace-nowrap px-3 py-2 text-sm text-muted-foreground">{review.pubYear}</td>
                  <td className="whitespace-nowrap px-3 py-2 text-sm text-muted-foreground">{review.dateRead}</td>
                  <td className="relative whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium flex justify-end">
                    <DropdownMenu modal={false}>
                      <DropdownMenuTrigger asChild>
                        <Button variant="outline">
                          <MoreVerticalIcon className="w-4 h-auto" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end" sideOffset={2}>
                        <>
                          <DropdownMenuItem onClick={() => handleEditReview(review)}>{t("common.update")}</DropdownMenuItem>
                          <DropdownMenuItem
                            onClick={() => handleDeleteReview(review)}
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

      {/* Create Review Dialog */}
      <CreateReviewDialog open={createDialog.isOpen} onOpenChange={createDialog.setOpen} onSuccess={fetchReviews} />

      {/* Edit Review Dialog */}
      <CreateReviewDialog open={editDialog.isOpen} onOpenChange={editDialog.setOpen} review={editingReview} onSuccess={fetchReviews} />
    </div>
  );
});

export default ReviewSection;
