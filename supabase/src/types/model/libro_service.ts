import { supabaseClient } from "@/store";
import { toItsError } from "./itserror";

export interface CountByGenre {
  genre: string;
  count: number;
}

export interface BookReview {
  id?: number;
  title: string;
  author: string;
  translator?: string | undefined;
  pages?: number | undefined;
  pubYear?: number | undefined;
  genre?: string | undefined;
  rating?: number | undefined;
  review?: string | undefined;
  dateRead?: string | undefined;
  createdAt?: Date | undefined;
}

function createBaseBookReview(): BookReview {
  return {
    title: "",
    author: "",
  };
}

export const BookReview: MessageFns<BookReview> = {
  create(base?: Partial<BookReview>): BookReview {
    return BookReview.fromPartial(base ?? {});
  },
  fromPartial(object: Partial<BookReview>): BookReview {
    const message = createBaseBookReview();
    message.title = object.title ?? "";
    message.author = object.author ?? "";
    message.translator = object.translator ?? "";
    if (object.pages) {
      message.pages = object.pages;
    }
    if (object.pubYear) {
      message.pubYear = object.pubYear;
    }
    message.genre = object.genre ?? "";
    if (object.rating) {
      message.rating = object.rating;
    }
    message.review = object.review ?? "";
    message.dateRead = object.dateRead ?? "";
    return message;
  },
};

export interface MessageFns<T> {
  create(base?: Partial<T>): T;
  fromPartial(object: Partial<T>): T;
}

const libroStore = (() => {
  const fetchBookReviews = async () => {
    let { data, error } = await supabaseClient
      .from("book_review")
      .select('*')
      .order('date_read', { ascending: false });
    return { data, error: error != null ? toItsError(error) : null };
  };

  const deleteBookReview = async (id: number) => {
    let { error } = await supabaseClient
      .from("book_review")
      .delete()
      .eq('id', id);
    return { error: error != null ? toItsError(error) : null }
  };

  const fetchStatsByGenre = async () => {
    let { data, error } = await supabaseClient
      .rpc('get_genre_all_stats');
    return { data, error };
  };

  return {
    fetchBookReviews,
    deleteBookReview,
    fetchStatsByGenre,
  };
})();

export { libroStore }
