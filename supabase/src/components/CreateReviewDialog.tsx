import { useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import useLoading from "@/hooks/useLoading";
import { BookReview, libroStore } from "@/types/model/libro_service";
import { toSnakeCase } from "@/utils/common";
import { useTranslate } from "@/utils/i18n";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  review?: BookReview;
  onSuccess?: () => void;
}

function CreateReviewDialog({ open, onOpenChange, review: initialReview, onSuccess }: Props) {
  const tableName = "book_review";
  const defaultDateRead = "1397-05-15";

  const t = useTranslate();
  const [review, setReview] = useState(BookReview.fromPartial({ ...initialReview }));
  const requestState = useLoading(false);
  const isCreating = !initialReview;

  useEffect(() => {
    if (initialReview) {
      setReview(BookReview.fromPartial(initialReview));
    } else {
      setReview(BookReview.fromPartial({}));
    }
  }, [initialReview]);

  const setPartialReview = (state: Partial<BookReview>) => {
    setReview({
      ...review,
      ...state,
    });
  };

  function isValidDateReadFormat(dateStr: string): boolean {
    if (!/^\d{4}-\d{2}-\d{2}$/.test(dateStr)) {
      console.error(dateStr + "success")
      return false;
    }

    const [year, month, day] = dateStr.split("-").map(Number);
    const date = new Date(year, month - 1, day);
    return (
      date.getFullYear() === year &&
      date.getMonth() === month - 1 &&
      date.getDate() === day
    );
  }

  const onPagesInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const pages = value === "" ? undefined : parseInt(value, 10) || undefined;

    setPartialReview({
      pages: pages,
    })
  };

  const onPubYearInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const pubYear = value === "" ? undefined : parseInt(value, 10) || undefined;

    setPartialReview({
      pubYear: pubYear,
    })
  };

  const onRatingInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const rating = value === "" ? undefined : parseInt(value, 10);

    setPartialReview({
      rating: Number.isNaN(rating) ? undefined : rating,
    })
  };

  const handleConfirm = async () => {
    if (isCreating && (!review.title || !review.author)) {
      toast.error("Title and author cannot be empty");
      return;
    }
    if (review.pages && review.pages > 999999) {
      toast.error("Pages should be between 1 and 999999");
      return;
    }
    if (review.pubYear && review.pubYear > 9999) {
      toast.error("Publication year should be between 1 and 9999");
      return;
    }
    if (review.rating && (review.rating < 0 || review.rating > 5)) {
      toast.error("Rating should be between 0 and 5");
      return;
    }
    if (review.dateRead !== undefined && !isValidDateReadFormat(review.dateRead)) {
      toast.error("Date read should be format like '" + defaultDateRead + "'");
      return;
    }

    try {
      requestState.setLoading();
      if (isCreating) {
        const { error } = await libroStore.insertReview(toSnakeCase(review))
        if (error != null) {
          throw error;
        }
        toast.success("Create book review successfully");
      } else {
        const updateReview: Record<string, any> = {};
        if (review.title !== initialReview?.title) {
          updateReview.title = review.title;
        }
        if (review.author !== initialReview?.author) {
          updateReview.author = review.author;
        }
        if (review.translator !== initialReview?.translator) {
          updateReview.translator = review.translator;
        }
        if (review.pages !== initialReview?.pages) {
          updateReview.pages = review.pages;
        }
        if (review.pubYear !== initialReview?.pubYear) {
          updateReview.pubYear = review.pubYear;
        }
        if (review.genre !== initialReview?.genre) {
          updateReview.genre = review.genre;
        }
        if (review.rating !== initialReview?.rating) {
          updateReview.rating = review.rating;
        }
        if (review.review !== initialReview?.review) {
          updateReview.review = review.review;
        }
        if (review.dateRead !== initialReview?.dateRead) {
          updateReview.dateRead = review.dateRead;
        }
        if (review.id !== undefined) {
          const { error } = await libroStore.updateReview(review.id, toSnakeCase(updateReview))
          if (error != null) {
            throw error;
          }
          toast.success("Update book review successfully");
        }
        
      }
      requestState.setFinish();
      onSuccess?.();
      onOpenChange(false);
      setReview(BookReview.fromPartial({}));
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
          <DialogTitle>{`${isCreating ? t("common.create") : t("common.edit")} ${t("common.review")}`}</DialogTitle>
          <DialogDescription>create or update a book review</DialogDescription>
        </DialogHeader>
        <div className="flex flex-col gap-4">
          <div className="grid gap-2">
            <Label htmlFor="title">{t("libro.title")}</Label>
            <Input
              id="title"
              type="text"
              placeholder={t("libro.title")}
              value={review.title}
              onChange={(e) =>
                setPartialReview({
                  title: e.target.value,
                })
              }
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="author">{t("libro.author")}</Label>
            <Input
              id="author"
              type="text"
              placeholder={t("libro.author")}
              value={review.author}
              onChange={(e) =>
                setPartialReview({
                  author: e.target.value,
                })
              }
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="translator">{t("libro.translator")}</Label>
            <Input
              id="translator"
              type="text"
              placeholder={t("libro.translator")}
              value={review.translator}
              onChange={(e) =>
                setPartialReview({
                  translator: e.target.value,
                })
              }
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="pages">{t("libro.pages")}</Label>
            <Input
              id="pages"
              type="text"
              placeholder={"0"}
              min={0}
              max={999999}
              value={review.pages ?? ''}
              onChange={onPagesInputChange}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="pub_year">{t("libro.pub-year")}</Label>
            <Input
              id="pub_year"
              type="text"
              placeholder={"1397"}
              min={0}
              max={9999}
              value={review.pubYear ?? ''}
              onChange={onPubYearInputChange}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="genre">{t("libro.genre")}</Label>
            <Input
              id="genre"
              type="text"
              placeholder={t("libro.genre")}
              value={review.genre}
              onChange={(e) =>
                setPartialReview({
                  genre: e.target.value,
                })
              }
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="rating">{t("libro.rating")}</Label>
            <Input
              id="rating"
              type="text"
              placeholder={"3"}
              min={0}
              max={5}
              value={review.rating ?? ''}
              onChange={onRatingInputChange}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="review">{t("libro.review")}</Label>
            <Input
              id="review"
              type="text"
              placeholder={t("libro.review")}
              value={review.review}
              onChange={(e) =>
                setPartialReview({
                  review: e.target.value,
                })
              }
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="date_read">{t("libro.date-read") + "(format: " + defaultDateRead + ")"}</Label>
            <Input
              id="date_read"
              type="text"
              placeholder={defaultDateRead}
              value={review.dateRead}
              onChange={(e) =>
                setPartialReview({
                  dateRead: e.target.value,
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

export default CreateReviewDialog;
