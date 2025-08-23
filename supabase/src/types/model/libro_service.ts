
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
