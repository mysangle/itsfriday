
-- user
INSERT INTO user (`username`, `role`, `email`, `nickname`, `password_hash`) VALUES ("skyblue", "USER", "its@friday.com skyblue", "skyblue", "$2a$10$gww237jWz4Gvvf7yKNhJz.5etRow4o2npbalhcOwXHjUSRv1QUXmi");

-- book
INSERT INTO book (`user_id`, `title`, `author`, `translator`, `pages`, `pub_year`, `genre`) VALUES (1, "질병 해방", "피터 아티아,빌 기퍼드", "이한음", 751, 2024, "health");
INSERT INTO book (`user_id`, `title`, `author`, `translator`, `pages`, `pub_year`, `genre`) VALUES (1, "협력의 진화", "로버트 액설로드", "이경식", 292, 2013, "science");

-- book review
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 1, "2024-05-01", 3, "Good");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 2, "2024-06-01", 4, "Excellent");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 1, "2024-07-01", 2, "Best");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 2, "2024-08-01", 4, "Diamond");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 1, "2024-09-01", 4, "Blue");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 2, "2024-10-01", 4, "Fantastic");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 1, "2024-11-01", 5, "Silver");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 2, "2024-12-01", 4, "Wonderful");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 1, "2025-01-01", 2, "Popular");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 2, "2025-02-01", 4, "Balance");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 1, "2025-03-01", 3, "For good lifestyle");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 2, "2025-04-01", 4, "One of the best book in my life");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 1, "2025-05-01", 4, "Best book in the world");
INSERT INTO book_review (`user_id`, `book_id`, `date_read`, `rating`, `review`) VALUES (1, 2, "2025-05-03", 1, "Good book");
