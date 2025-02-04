-- +goose Up


INSERT INTO public.users
(id, first_name, last_name, birthdate, biography, city, username, password_hash)
VALUES('d86c3ba3-0628-4716-88b8-a51ec02bc603'::uuid, 'John', 'Doe', '2000-01-01', 'Blah-blah-blah', 'Silent Hill', 'johndoe@gmail.com', '482c811da5d5b4bc6d497ffa98491e38');
INSERT INTO public.users
(id, first_name, last_name, birthdate, biography, city, username, password_hash)
VALUES('b3dd0beb-26a1-4437-a7bd-3a57e40bb4e3'::uuid, 'Maria', 'Pupkina', '2000-02-03', 'Blah-blah-blah', 'Silent Hill', 'masha@gmail.com', '482c811da5d5b4bc6d497ffa98491e38');

INSERT INTO public.users_friends
(id, friend_id)
VALUES('b3dd0beb-26a1-4437-a7bd-3a57e40bb4e3'::uuid, 'd86c3ba3-0628-4716-88b8-a51ec02bc603'::uuid);
INSERT INTO public.users_friends
(id, friend_id)
VALUES('d86c3ba3-0628-4716-88b8-a51ec02bc603'::uuid, 'b3dd0beb-26a1-4437-a7bd-3a57e40bb4e3'::uuid);

-- +goose Down

DELETE FROM public.users
WHERE id IN ('d86c3ba3-0628-4716-88b8-a51ec02bc603','b3dd0beb-26a1-4437-a7bd-3a57e40bb4e3' );

DELETE FROM users_friends 
WHERE id IN ('d86c3ba3-0628-4716-88b8-a51ec02bc603','b3dd0beb-26a1-4437-a7bd-3a57e40bb4e3' ) 
    AND friend_id IN ('d86c3ba3-0628-4716-88b8-a51ec02bc603','b3dd0beb-26a1-4437-a7bd-3a57e40bb4e3' );
