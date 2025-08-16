
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
    message.id = object.id ?? 0;
    message.title = object.title ?? "";
    message.author = object.author ?? "";
    message.translator = object.translator ?? "";
    message.pages = object.pages ?? 0;
    message.pubYear = object.pubYear ?? 0;
    message.genre = object.genre ?? "";
    message.rating = object.rating ?? 0;
    message.review = object.review ?? "";
    message.dateRead = object.dateRead ?? "";
    message.createdAt = object.createdAt;
    return message;
  },
};

export interface MessageFns<T> {
  create(base?: Partial<T>): T;
  fromPartial(object: Partial<T>): T;
}
