-- +goose Up
insert into access_list (app_id, method, value)
values (1, 'admin/auth/login', true)
on conflict (app_id, method) do update set value = true;

-- +goose Down

