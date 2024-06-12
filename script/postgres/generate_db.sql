-- person definition

-- Drop table

drop table if exists person_role;
drop table if exists film;
drop table if exists studio;
drop table if exists person;

CREATE TABLE person
(
    id         serial4 NOT NULL,
    "name"     varchar NOT NULL,
    birth_date date    NOT NULL,
    CONSTRAINT person_pk PRIMARY KEY (id)
);


-- rating definition

CREATE TABLE rating
(
    id     int4    NOT NULL,
    "name" varchar NOT NULL,
    CONSTRAINT rating_pk PRIMARY KEY (id)
);


-- studio definition

CREATE TABLE studio
(
    id     serial4 NOT NULL,
    "name" varchar NOT NULL,
    CONSTRAINT studio_pk PRIMARY KEY (id)
);


-- film definition

CREATE TABLE film
(
    id        serial4        NOT NULL,
    "name"    varchar        NOT NULL,
    "year"    int4           NOT NULL,
    rating_id int4           NOT NULL,
    gross     int4 DEFAULT 0 NOT NULL,
    studio_id int4           NOT NULL,
    CONSTRAINT film_check CHECK ((year >= 1800)),
    CONSTRAINT film_pk PRIMARY KEY (id),
    CONSTRAINT film_unique UNIQUE (year, name),
    CONSTRAINT film_rating_fk FOREIGN KEY (rating_id) REFERENCES rating (id) ON DELETE SET DEFAULT ON UPDATE CASCADE,
    CONSTRAINT film_studio_fk FOREIGN KEY (studio_id) REFERENCES studio (id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- person_role definition

CREATE TABLE person_role
(
    id        serial4 NOT NULL,
    person_id int4    NOT NULL,
    film_id   int4    NOT NULL,
    "role"    varchar NOT NULL,
    CONSTRAINT person_role_pk PRIMARY KEY (id),
    CONSTRAINT person_role_film_fk FOREIGN KEY (film_id) REFERENCES film (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT person_role_person_fk FOREIGN KEY (person_id) REFERENCES person (id) ON DELETE CASCADE ON UPDATE CASCADE
);

BEGIN;
INSERT INTO rating (id, name) VALUES (1, 'PG-10');
INSERT INTO rating (id, name) VALUES (2, 'PG-13');
INSERT INTO rating (id, name) VALUES (3, 'PG-18');
COMMIT;

BEGIN;
INSERT INTO studio (name) VALUES ('Universal Pictures'); -- 1
INSERT INTO studio (name) VALUES ('Paramount Pictures'); -- 2
INSERT INTO studio (name) VALUES ('Warner Bros. Pictures'); -- 3
INSERT INTO studio (name) VALUES ('Walt Disney Pictures'); -- 4
INSERT INTO studio (name) VALUES ('Columbia Pictures'); -- 5
COMMIT;

BEGIN;
INSERT INTO film (name, year, rating_id, gross, studio_id) VALUES ('The Grand Budapest Hotel', 2014, 2, 800, 1); -- 1
INSERT INTO film (name, year, rating_id, gross, studio_id) VALUES ('Django Unchained', 2012, 3, 2000, 2); -- 2
INSERT INTO film (name, year, rating_id, gross, studio_id) VALUES ('Once Upon a Time in... Hollywood', 2019, 3, 1500, 2); -- 3
INSERT INTO film (name, year, rating_id, gross, studio_id) VALUES ('Sen to Chihiro no kamikakushi', 2001, 1, 3500, 4); -- 4

INSERT INTO person (name, birth_date) VALUES('Wes Anderson', 'May 01, 1969'); -- 0
INSERT INTO person (name, birth_date) VALUES('Ralph Fiennes', 'December 22, 1962'); -- 1
INSERT INTO person (name, birth_date) VALUES('Quentin Tarantino', 'March 27, 1963'); -- 2
INSERT INTO person (name, birth_date) VALUES('Leonardo DiCaprio', 'March 27, 1963'); -- 3
INSERT INTO person (name, birth_date) VALUES('Hayao Miyazaki', 'January 5, 1941'); -- 4

INSERT INTO person_role (person_id, film_id, role) VALUES (1, 1, 'director');
INSERT INTO person_role (person_id, film_id, role) VALUES (2, 1, 'actor');

INSERT INTO person_role (person_id, film_id, role) VALUES (3, 2, 'director');
INSERT INTO person_role (person_id, film_id, role) VALUES (4, 2, 'actor');

INSERT INTO person_role (person_id, film_id, role) VALUES (3, 3, 'director');
INSERT INTO person_role (person_id, film_id, role) VALUES (4, 3, 'actor');

INSERT INTO person_role (person_id, film_id, role) VALUES (5, 4, 'director');
END;

SELECT * FROM film;
SELECT * FROM person;
SELECT * FROM studio;